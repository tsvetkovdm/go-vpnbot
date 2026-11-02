package payments_transport_tg

import (
	"context"
	"errors"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
	"go.uber.org/zap"
)

func (h *PaymentsTGHandler) MenuClaimFreePayment(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuClaimFreePayment handler")

	messageID := update.CallbackQuery.Message.MessageID
	telegramID := core_tg_request.FromUser(update).ID
	backMarkup := noPaymentMarkup()

	priceID, err := parseIDArg(update.CallbackQuery.Data)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &backMarkup)
		return
	}

	resolved, err := h.paymentsService.ClaimFreePayment(ctx, telegramID, priceID)
	switch {
	case err == nil:
		paidUntil := resolved.Subscription.PaidUntil(time.Local).Format("02.01.2006")
		markup := tgbotapi.NewInlineKeyboardMarkup(accountRow())
		responseHandler.EditResponse(
			fmt.Sprintf("🎁 Готово! Бесплатный доступ открыт до %s включительно.", paidUntil),
			&markup,
		)
		return

	case errors.Is(err, domain.ErrPaymentFreeUnavailable):
		log.Info("free price is unavailable for user", zap.Error(err))
		responseHandler.EditResponse("🎁 Бесплатный доступ — только для новых пользователей.", &backMarkup)
		return

	case errors.Is(err, domain.ErrPaymentPlanUnavailable):
		log.Info("plan price unavailable", zap.Error(err))
		responseHandler.EditResponse("Этот вариант больше недоступен.", &backMarkup)
		return

	case errors.Is(err, domain.ErrPaymentUserBlocked):
		log.Info("blocked user tried to claim free price", zap.Error(err))
		homeMarkup := tgbotapi.NewInlineKeyboardMarkup(accountRow())
		responseHandler.EditResponse("🚫 Доступ ограничен.\nЕсли это ошибка — напишите "+h.adminContact, &homeMarkup)
		return

	case errors.Is(err, core_errors.ErrNotFound):
		responseHandler.EditResponse("Нажмите /start, чтобы начать.", nil)
		return

	default:
		responseHandler.EditErrorResponse(messageID, err, "failed to claim free payment", &backMarkup)
		return
	}
}
