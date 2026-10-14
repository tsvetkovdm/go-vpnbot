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

			session, status := store.Get(chatID)

			switch status {
			case core_tg_dialog.NotFound:
				next(ctx, update)
				return
			case core_tg_dialog.Expired:
				core_tg_response.NewTGResponseHandler(log, update, bot).
					TextResponse("Сессия истекла, начните заново")
				next(ctx, update)
				return
			}

			if update.Message == nil || update.Message.IsCommand() {
				store.Delete(chatID)
				if update.Message != nil && update.Message.IsCommand() {
					core_tg_response.NewTGResponseHandler(log, update, bot).
						EditMessage(session.MessageID, "Создание тарифа отменено", nil)
				}
				next(ctx, update)
				return
			}

			delMsg := tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID)
			_, err := bot.Send(delMsg)
			if err != nil {
				log.Debug("dialog: failed to delete user message", zap.Error(err))
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
			responseHandler := core_tg_response.NewTGResponseHandler(log, update, bot)

			switch result.Kind {
			case core_tg_dialog.Retry:
				log.Debug("dialog: retry",
					zap.String("step_id", string(session.Step)),
					zap.String("reason", result.Message),
				)
				store.Set(chatID, session)
				responseHandler.EditMessage(session.MessageID, result.Message, nil)

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
				responseHandler.EditMessage(session.MessageID, enterText, nil)

			case core_tg_dialog.Complete:
				log.Debug("dialog: complete",
					zap.String("flow_id", session.FlowID),
					zap.String("step_id", string(session.Step)),
				)
				store.Delete(chatID)
				responseHandler.EditMessage(session.MessageID, result.Message, nil)

			case core_tg_dialog.Abort:
				store.Delete(chatID)
				responseHandler.EditErrorResponse(session.MessageID, result.Err, result.LogMsg, nil)
			}
		}
	}
}
