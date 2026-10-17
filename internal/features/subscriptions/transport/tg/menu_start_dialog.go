package subscriptions_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *SubscriptionsTGHandler) MenuStartDialog(ctx context.Context, update tgbotapi.Update) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)
	log.Debug("invoke MenuStartDialog handler")

	chatID := core_tg_request.ChatID(update)
	messageID := update.CallbackQuery.Message.MessageID
	data := update.CallbackQuery.Data

	target, err := parseDialogCallback(data)
	if err != nil {
		adminMenu := adminMenuMarkup()
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &adminMenu)
		return
	}

	flow := h.extendDaysFlow
	if core_tg_callback.Segments(data)[1] == "date" {
		flow = h.setDateFlow
	}

	step := flow.Steps[flow.Start]
	markup := core_tg_dialog.StepMarkup(step)
	responseHandler.EditResponse(step.Enter, &markup)

	session := core_tg_dialog.NewSession(flow.ID, messageID, flow.Start)
	target.saveTo(&session)
	h.dialogStore.Set(chatID, session)
}
