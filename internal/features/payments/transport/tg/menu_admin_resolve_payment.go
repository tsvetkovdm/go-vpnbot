package payments_transport_tg

import (
	"context"
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
	"go.uber.org/zap"
)

func (h *PaymentsTGHandler) MenuAdminResolvePayment(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuAdminResolvePayment handler")

	messageID := update.CallbackQuery.Message.MessageID
	backMarkup := adminPaymentsMarkup(nil)

	paymentID, err := parseIDArg(update.CallbackQuery.Data)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &backMarkup)
		return
	}

	adminTelegramID := core_tg_request.FromUser(update).ID
	confirm := callbackAction(update.CallbackQuery.Data) == "aconfirm"

	var (
		resolved domain.ResolvedPayment
		answer   string
	)
	if confirm {
		resolved, err = h.paymentsService.ConfirmPayment(ctx, paymentID, adminTelegramID)
	} else {
		resolved, err = h.paymentsService.RejectPayment(ctx, paymentID, adminTelegramID)
	}

	switch {
	case errors.Is(err, domain.ErrPaymentAlreadyResolved):
		log.Info("payment already resolved", zap.Error(err))
		core_tg_callback.SetAnswer(ctx, "ℹ️ Оплата уже обработана")
		h.showAdminPayments(ctx, responseHandler, update)
		return

	case errors.Is(err, core_errors.ErrConflict):
		log.Info("subscription changed concurrently", zap.Error(err))
		core_tg_callback.SetAnswer(ctx, "⚠️ Данные изменились, нажмите ещё раз")
		return

	case err != nil:
		log.Error("failed to resolve payment", zap.Error(err))
		core_tg_callback.SetAnswer(ctx, "❌ Ошибка, попробуйте позже")
		return
	}

	if confirm {
		answer = h.sendConfirmedToUser(ctx, resolved)
	} else {
		answer = h.sendRejectedToUser(ctx, resolved)
	}

	core_tg_callback.SetAnswer(ctx, answer)
	h.showAdminPayments(ctx, responseHandler, update)
}
