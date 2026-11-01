package payments_transport_tg

import (
	"context"
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_html "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/html"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
	"go.uber.org/zap"
)

func (h *PaymentsTGHandler) MenuCreatePayment(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuCreatePayment handler")

	messageID := update.CallbackQuery.Message.MessageID
	telegramID := core_tg_request.FromUser(update).ID
	homeRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("👤 Личный кабинет", "users:home"),
	)
	homeMarkup := tgbotapi.NewInlineKeyboardMarkup(homeRow)
	backMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ К тарифам", "plans:public"),
		),
		homeRow,
	)

	priceID, err := parseIDArg(update.CallbackQuery.Data)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &backMarkup)
		return
	}

	payment, err := h.paymentsService.CreatePayment(ctx, telegramID, priceID)
	switch {
	case err == nil:
		markup := paymentMarkup(payment)
		responseHandler.EditHTMLResponse(h.formatPayment(payment), &markup)
		return

	case errors.Is(err, domain.ErrPaymentUnresolvedExists):
		log.Info("user already has unresolved payment", zap.Error(err))

		existing, getErr := h.paymentsService.GetUnresolvedPayment(ctx, telegramID)
		if getErr != nil {
			responseHandler.EditErrorResponse(messageID, getErr, "failed to get unresolved payment", &backMarkup)
			return
		}

		text := core_tg_html.Sprintf(
			"⚠️ У вас уже есть неоплаченный счёт.\n\n%s",
			h.formatPayment(existing),
		)
		markup := paymentMarkup(existing)
		responseHandler.EditHTMLResponse(text, &markup)
		return

	case errors.Is(err, domain.ErrPaymentPlanUnavailable):
		log.Info("plan price unavailable", zap.Error(err))
		responseHandler.EditResponse("Этот вариант больше недоступен.", &backMarkup)
		return

	case errors.Is(err, domain.ErrPaymentUserBlocked):
		log.Info("blocked user tried to pay", zap.Error(err))
		responseHandler.EditResponse("🚫 Доступ ограничен.\nЕсли это ошибка — напишите "+h.adminContact, &homeMarkup)
		return

	case errors.Is(err, core_errors.ErrNotFound):
		responseHandler.EditResponse("Нажмите /start, чтобы начать.", nil)
		return

	default:
		responseHandler.EditErrorResponse(messageID, err, "failed to create payment", &backMarkup)
		return
	}
}
