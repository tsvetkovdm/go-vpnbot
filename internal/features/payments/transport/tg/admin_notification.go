package payments_transport_tg

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
	"go.uber.org/zap"
)

func (h *PaymentsTGHandler) MenuDismissNotification(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuDismissNotification handler")

	h.closeNotification(responseHandler, update, "ℹ️ Оплата уже обработана")
}

func (h *PaymentsTGHandler) closeNotification(
	responseHandler *core_tg_response.TGResponseHandler,
	update tgbotapi.Update,
	fallbackText string,
) {
	chatID := core_tg_request.ChatID(update)
	messageID := update.CallbackQuery.Message.MessageID

	if _, err := h.bot.Request(tgbotapi.NewDeleteMessage(chatID, messageID)); err == nil {
		return
	}

	responseHandler.EditResponse(fallbackText, nil)
}

func (h *PaymentsTGHandler) showAlreadyResolved(
	ctx context.Context,
	responseHandler *core_tg_response.TGResponseHandler,
	paymentID int64,
) {
	log := core_logger.FromContext(ctx)

	okMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ок", "pay:dismiss"),
		),
	)

	payment, err := h.paymentsService.GetPaymentByID(ctx, paymentID)
	if err != nil {
		log.Error("failed to get payment", zap.Error(err))
		responseHandler.EditResponse("ℹ️ Эта оплата уже обработана.", &okMarkup)
		return
	}

	responseHandler.EditResponse(
		fmt.Sprintf(
			"ℹ️ Оплата с кодом %d уже обработана: %s.",
			payment.CommentCode, paymentStatusText(payment.Status),
		),
		&okMarkup,
	)
}

func paymentStatusText(status domain.PaymentStatus) string {
	switch status {
	case domain.PaymentStatusSuccessful:
		return "подтверждена"
	case domain.PaymentStatusRejected:
		return "отклонена"
	case domain.PaymentStatusCancelled:
		return "отменена пользователем"
	case domain.PaymentStatusExpired:
		return "истекла"
	default:
		return string(status)
	}
}
