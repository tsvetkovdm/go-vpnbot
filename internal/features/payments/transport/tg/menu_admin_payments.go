package payments_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *PaymentsTGHandler) MenuAdminPayments(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuAdminPayments handler")

	h.showAdminPayments(ctx, responseHandler, update)
}

func (h *PaymentsTGHandler) showAdminPayments(
	ctx context.Context,
	responseHandler *core_tg_response.TGResponseHandler,
	update tgbotapi.Update,
) {
	messageID := update.CallbackQuery.Message.MessageID

	overviews, total, err := h.paymentsService.ListPendingPayments(ctx, adminPaymentsLimit)
	if err != nil {
		markup := adminPaymentsMarkup(nil)
		responseHandler.EditErrorResponse(messageID, err, "failed to list pending payments", &markup)
		return
	}

	markup := adminPaymentsMarkup(overviews)
	responseHandler.EditResponse(adminPaymentsText(len(overviews), total), &markup)
}
