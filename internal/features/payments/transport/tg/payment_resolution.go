package payments_transport_tg

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (h *PaymentsTGHandler) sendConfirmedToUser(
	ctx context.Context,
	resolved domain.ResolvedPayment,
) string {
	paidUntil := resolved.Subscription.PaidUntil(time.Local).Format("02.01.2006")
	p := resolved.Payment

	h.screens.Replace(
		ctx,
		resolved.UserTelegramID,
		fmt.Sprintf(
			"✅ Оплата подтверждена!\n\n%s, %s\nДоступ до %s включительно.",
			p.PlanTitle, formatDuration(p.Duration), paidUntil,
		),
		tgbotapi.NewInlineKeyboardMarkup(accountRow()),
	)

	return fmt.Sprintf("✅ Подтверждено: код %d, доступ до %s", p.CommentCode, paidUntil)
}

func (h *PaymentsTGHandler) sendRejectedToUser(
	ctx context.Context,
	resolved domain.ResolvedPayment,
) string {
	p := resolved.Payment

	h.screens.Replace(
		ctx,
		resolved.UserTelegramID,
		fmt.Sprintf(
			"❌ Перевод не найден\n\n"+
				"Мы не нашли перевод с комментарием %d на %s.\n"+
				"Если вы уже перевели деньги — напишите %s и назовите этот код.",
			p.CommentCode, formatAmount(p.Price), h.adminContact,
		),
		noPaymentMarkup(),
	)

	return fmt.Sprintf("❌ Отклонено: код %d", p.CommentCode)
}
