package payments_transport_tg

import (
	"context"
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
	"go.uber.org/zap"
)

func (h *PaymentsTGHandler) MenuSubmitPayment(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuSubmitPayment handler")

	messageID := update.CallbackQuery.Message.MessageID
	from := core_tg_request.FromUser(update)
	errorMarkup := noPaymentMarkup()

	paymentID, err := parseIDArg(update.CallbackQuery.Data)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &errorMarkup)
		return
	}

	payment, err := h.paymentsService.SubmitPayment(ctx, from.ID, paymentID)
	if errors.Is(err, core_errors.ErrConflict) {
		log.Info("payment is not submittable", zap.Error(err))
		h.showCurrentPayment(ctx, responseHandler, update)
		return
	}
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to submit payment", &errorMarkup)
		return
	}

	h.notifyAdmin(ctx, payment, from)

	markup := paymentMarkup(payment)
	responseHandler.EditHTMLResponse(h.formatPayment(payment), &markup)
}
