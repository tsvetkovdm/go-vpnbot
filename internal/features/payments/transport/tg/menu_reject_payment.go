package payments_transport_tg

import (
	"context"
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
	"go.uber.org/zap"
)

func (h *PaymentsTGHandler) MenuRejectPayment(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuRejectPayment handler")

	paymentID, err := parseIDArg(update.CallbackQuery.Data)
	if err != nil {
		log.Error("failed to parse callback data", zap.Error(err))
		core_tg_callback.SetAnswer(ctx, "❌ Не удалось разобрать кнопку")
		return
	}

	adminTelegramID := core_tg_request.FromUser(update).ID

	resolved, err := h.paymentsService.RejectPayment(ctx, paymentID, adminTelegramID)
	if errors.Is(err, domain.ErrPaymentAlreadyResolved) {
		log.Info("payment already resolved", zap.Error(err))
		h.showAlreadyResolved(ctx, responseHandler, paymentID)
		return
	}
	if err != nil {
		log.Error("failed to reject payment", zap.Error(err))
		core_tg_callback.SetAnswer(ctx, "❌ Ошибка, попробуйте позже")
		return
	}

	h.closeNotification(responseHandler, update, "❌ Отклонено")
	core_tg_callback.SetAnswer(ctx, h.sendRejectedToUser(ctx, resolved))
}
