package payments_transport_tg

import (
	"context"
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *PaymentsTGHandler) MenuAdminPayment(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuAdminPayment handler")

	messageID := update.CallbackQuery.Message.MessageID
	backMarkup := adminPaymentsMarkup(nil)

	paymentID, err := parseIDArg(update.CallbackQuery.Data)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &backMarkup)
		return
	}

	overview, err := h.paymentsService.GetPaymentOverview(ctx, paymentID)
	if errors.Is(err, core_errors.ErrNotFound) {
		core_tg_callback.SetAnswer(ctx, "Оплата не найдена")
		h.showAdminPayments(ctx, responseHandler, update)
		return
	}
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to get payment overview", &backMarkup)
		return
	}

	markup := adminPaymentCardMarkup(overview.Payment)
	responseHandler.EditHTMLResponse(formatAdminPaymentCard(overview), &markup)
}
