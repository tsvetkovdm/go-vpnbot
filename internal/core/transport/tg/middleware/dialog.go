package core_tg_middleware

import (
	"context"
	"strings"

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
			chatID := core_tg_request.ChatID(update)

			session, status := store.Get(chatID)
			if status == core_tg_dialog.NotFound {
				if update.CallbackQuery != nil &&
					strings.HasPrefix(update.CallbackQuery.Data, core_tg_dialog.CallbackPrefix) {
					log := core_logger.FromContext(ctx)
					core_tg_response.NewTGResponseHandler(log, update, bot).
						EditResponse("Это действие больше не активно.", nil)
					return
				}
				next(ctx, update)
				return
			}

			log := core_logger.FromContext(ctx)
			responseHandler := core_tg_response.NewTGResponseHandler(log, update, bot)

			flow, ok := flows[session.FlowID]
			if !ok {
				log.Error("dialog: unknown flow id in session", zap.String("flow_id", session.FlowID))
				store.Delete(chatID)
				next(ctx, update)
				return
			}

			if status == core_tg_dialog.Expired {
				responseHandler.EditMessage(session.MessageID, "Время ожидания истекло. Начните заново.", flow.ExitMarkup)
				next(ctx, update)
				return
			}

			input, ok := dialogInput(update)
			if !ok {
				store.Delete(chatID)
				responseHandler.EditMessage(session.MessageID, flow.CancelText, flow.ExitMarkup)
				next(ctx, update)
				return
			}

			if update.CallbackQuery != nil && input == core_tg_dialog.CancelValue {
				store.Delete(chatID)
				responseHandler.EditMessage(session.MessageID, flow.CancelText, flow.ExitMarkup)
				return
			}

			if update.Message != nil {
				deleteMsg := tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID)
				if _, err := bot.Request(deleteMsg); err != nil {
					log.Debug("dialog: failed to delete user message", zap.Error(err))
				}
			}

			step, ok := flow.Steps[session.Step]
			if !ok {
				log.Error("dialog: unknown step id in session", zap.String("step_id", string(session.Step)))
				store.Delete(chatID)
				responseHandler.EditMessage(session.MessageID, flow.CancelText, flow.ExitMarkup)
				return
			}

			result := step.OnInput(ctx, &session, input)

			switch result.Kind {
			case core_tg_dialog.Retry:
				store.Set(chatID, session)
				markup := core_tg_dialog.StepMarkup(step)
				responseHandler.EditMessage(session.MessageID, result.Message+"\n\n"+step.Enter, &markup)

			case core_tg_dialog.Advance:
				nextStep, ok := flow.Steps[result.Next]
				if !ok {
					log.Error("dialog: advance to unknown step", zap.String("step_id", string(result.Next)))
					store.Delete(chatID)
					responseHandler.EditMessage(session.MessageID, flow.CancelText, flow.ExitMarkup)
					return
				}

				session.Step = result.Next
				store.Set(chatID, session)

				enterText := nextStep.Enter
				if result.Message != "" {
					enterText = result.Message
				}
				markup := core_tg_dialog.StepMarkup(nextStep)
				responseHandler.EditMessage(session.MessageID, enterText, &markup)

			case core_tg_dialog.Complete:
				store.Delete(chatID)
				markup := flow.ExitMarkup
				if result.Markup != nil {
					markup = result.Markup
				}
				responseHandler.EditMessage(session.MessageID, result.Message, markup)

			case core_tg_dialog.Abort:
				store.Delete(chatID)
				responseHandler.EditErrorResponse(session.MessageID, result.Err, result.LogMsg, flow.ExitMarkup)
			}
		}
	}
}

func dialogInput(update tgbotapi.Update) (string, bool) {
	switch {
	case update.Message != nil && !update.Message.IsCommand():
		return update.Message.Text, true
	case update.CallbackQuery != nil &&
		strings.HasPrefix(update.CallbackQuery.Data, core_tg_dialog.CallbackPrefix):
		return strings.TrimPrefix(update.CallbackQuery.Data, core_tg_dialog.CallbackPrefix), true
	}

	return "", false
}
