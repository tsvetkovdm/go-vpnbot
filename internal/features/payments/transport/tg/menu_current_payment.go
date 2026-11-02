package payments_transport_tg

import (
	"context"
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *PaymentsTGHandler) MenuCurrentPayment(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuCurrentPayment handler")

	h.showCurrentPayment(ctx, responseHandler, update)
}

func (h *PaymentsTGHandler) showCurrentPayment(
	ctx context.Context,
	responseHandler *core_tg_response.TGResponseHandler,
	update tgbotapi.Update,
) {
	messageID := update.CallbackQuery.Message.MessageID
	telegramID := core_tg_request.FromUser(update).ID

	payment, err := h.paymentsService.GetUnresolvedPayment(ctx, telegramID)
	if errors.Is(err, core_errors.ErrNotFound) {
		markup := noPaymentMarkup()
		responseHandler.EditResponse("Эта оплата уже завершена или отменена.", &markup)
		return
	}
	if err != nil {
		markup := noPaymentMarkup()
		responseHandler.EditErrorResponse(messageID, err, "failed to get unresolved payment", &markup)
		return
	}

	markup := paymentMarkup(payment)
	responseHandler.EditHTMLResponse(h.formatPayment(payment), &markup)
}
