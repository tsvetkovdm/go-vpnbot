package payments_transport_tg

import (
	"context"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_html "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/html"
	"go.uber.org/zap"
)

func (h *PaymentsTGHandler) notifyAdmin(
	ctx context.Context,
	payment domain.Payment,
	from *tgbotapi.User,
) {
	log := core_logger.FromContext(ctx)

	msg := tgbotapi.NewMessage(h.adminChatID, string(formatAdminNotification(payment, from)))
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = adminNotificationMarkup(payment)

	if _, err := h.bot.Send(msg); err != nil {
		log.Error(
			"failed to notify admin about payment",
			zap.Int64("payment_id", payment.ID),
			zap.Error(err),
		)
	}
}

func formatAdminNotification(p domain.Payment, from *tgbotapi.User) core_tg_html.HTML {
	payer := core_tg_html.Escape(from.FirstName)
	if from.UserName != "" {
		payer = core_tg_html.Sprintf("%s (@%s)", payer, core_tg_html.Escape(from.UserName))
	}

	return core_tg_html.Sprintf(
		"🔔 Новая оплата\n\n"+
			"Код: %s\n"+
			"Сумма: %s\n"+
			"Тариф: %s, %s\n"+
			"От: %s, ID %s\n\n"+
			"Проверьте поступление и подтвердите или отклоните.",
		core_tg_html.Code(strconv.Itoa(p.CommentCode)),
		core_tg_html.Bold(formatAmount(p.Price)),
		core_tg_html.Escape(p.PlanTitle),
		core_tg_html.Escape(formatDuration(p.Duration)),
		payer,
		core_tg_html.Code(strconv.FormatInt(from.ID, 10)),
	)
}

func adminNotificationMarkup(p domain.Payment) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Подтвердить", paymentCallback("confirm", p.ID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отклонить", paymentCallback("reject", p.ID)),
		),
	)
}
