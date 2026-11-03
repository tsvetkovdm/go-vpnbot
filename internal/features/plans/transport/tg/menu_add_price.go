package plans_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *PlansTGHandler) MenuAddPrice(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuAddPrice handler")

	chatID := core_tg_request.ChatID(update)
	messageID := update.CallbackQuery.Message.MessageID

	planID, filter, page, err := parseCallbackData(update.CallbackQuery.Data)
	if err != nil {
		markup := plansMenuMarkup()
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &markup)
		return
	}

	flow := h.addPriceFlow
	step := flow.Steps[flow.Start]
	markup := core_tg_dialog.StepMarkup(step)
	responseHandler.EditResponse(step.Enter, &markup)

	session := core_tg_dialog.NewSession(flow.ID, messageID, flow.Start)
	planTarget{planID: planID, filter: filter, page: page}.saveTo(&session)
	h.dialogStore.Set(chatID, session)
}
