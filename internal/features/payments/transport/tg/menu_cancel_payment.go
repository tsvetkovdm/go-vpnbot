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

func (h *PaymentsTGHandler) MenuCancelPayment(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuCancelPayment handler")

	messageID := update.CallbackQuery.Message.MessageID
	telegramID := core_tg_request.FromUser(update).ID
	doneMarkup := noPaymentMarkup()

	paymentID, err := parseIDArg(update.CallbackQuery.Data)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &doneMarkup)
		return
	}

	_, err = h.paymentsService.CancelPayment(ctx, telegramID, paymentID)
	if errors.Is(err, core_errors.ErrConflict) {
		log.Info("payment is not cancellable", zap.Error(err))

		markup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("💳 Продолжить оплату", "pay:current"),
			),
			accountRow(),
		)
		responseHandler.EditResponse("⚠️ Отменить уже нельзя — оплата изменилась.", &markup)
		return
	}
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to cancel payment", &doneMarkup)
		return
	}

	responseHandler.EditResponse("Оплата отменена. Можно выбрать другой тариф.", &doneMarkup)
}
