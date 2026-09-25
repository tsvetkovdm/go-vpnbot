package core_tg_middleware

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
	core_tg_handler "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/handler"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
	"go.uber.org/zap"
)

func Dialog(
	bot *tgbotapi.BotAPI,
	store *core_tg_dialog.Store,
	flows map[string]core_tg_dialog.Flow,
) Middleware {
	return func(next core_tg_handler.Handler) core_tg_handler.Handler {
		return func(ctx context.Context, update tgbotapi.Update) {
			log := core_logger.FromContext(ctx)
			chatID := core_tg_request.ChatID(update)
			responseHandler := core_tg_response.NewTGResponseHandler(log, update, bot)

			session, status := store.Get(chatID)
			switch status {
			case core_tg_dialog.NotFound:
				next(ctx, update)
				return
			case core_tg_dialog.Expired:
				responseHandler.TextResponse("Сессия истекла, начните заново")
				next(ctx, update)
				return
			}

			if update.Message == nil || update.Message.IsCommand() {
				store.Delete(chatID)
				if update.Message != nil && update.Message.IsCommand() {
					responseHandler.TextResponse("Создание тарифа отменено")
				}
				next(ctx, update)
				return
			}

			flow, ok := flows[session.FlowID]
			if !ok {
				log.Error("dialog: unknown flow id in session", zap.String("flow_id", session.FlowID))
				store.Delete(chatID)
				return
			}

			step, ok := flow.Steps[session.Step]
			if !ok {
				log.Error("dialog: unknown step id in session", zap.String("step_id", string(session.Step)))
				store.Delete(chatID)
				return
			}

			result := step.OnInput(ctx, &session, update)

			switch result.Kind {
			case core_tg_dialog.Retry:
				log.Debug("dialog: retry",
					zap.String("step_id", string(session.Step)),
					zap.String("reason", result.Message),
				)
				store.Set(chatID, session)
				responseHandler.TextResponse(result.Message)

			case core_tg_dialog.Advance:
				nextStep, ok := flow.Steps[result.Next]
				if !ok {
					log.Error("dialog: advance to unknown step", zap.String("step_id", string(result.Next)))
					store.Delete(chatID)
					return
				}

				log.Debug("dialog: advance",
					zap.String("from_step", string(session.Step)),
					zap.String("to_step", string(result.Next)),
				)

				session.Step = result.Next
				store.Set(chatID, session)

				enterText := nextStep.Enter
				if result.Message != "" {
					enterText = result.Message
				}
				responseHandler.TextResponse(enterText)

			case core_tg_dialog.Complete:
				log.Debug("dialog: complete",
					zap.String("flow_id", session.FlowID),
					zap.String("step_id", string(session.Step)),
				)
				store.Delete(chatID)
				responseHandler.TextResponse(result.Message)
			case core_tg_dialog.Abort:
				store.Delete(chatID)
				responseHandler.ErrorResponse(result.Err, result.LogMsg)
			}
		}
	}
}
