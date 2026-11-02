package payments_transport_tg

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_html "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/html"
)

const adminPaymentsLimit = 20

func adminPaymentsText(shown, total int) string {
	if total == 0 {
		return "💰 Оплаты на проверке\n\nНовых оплат нет."
	}

	text := fmt.Sprintf("💰 Оплаты на проверке: %d\n\nСначала самые давние.", total)
	if total > shown {
		text += fmt.Sprintf("\nПоказаны первые %d.", shown)
	}
	return text
}

func adminPaymentsMarkup(overviews []domain.PaymentOverview) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(overviews)+2)

	for _, o := range overviews {
		label := fmt.Sprintf(
			"%d · %s · %s",
			o.Payment.CommentCode, formatAmount(o.Payment.Price), o.UserFirstName,
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, paymentCallback("aview", o.Payment.ID)),
		))
	}

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить", "pay:admin"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Админ-меню", "admins:menu"),
		),
	)

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func formatAdminPaymentCard(o domain.PaymentOverview) core_tg_html.HTML {
	username := ""
	if o.UserUsername != nil {
		username = *o.UserUsername
	}

	body := formatPaymentForAdmin(o.Payment, o.UserTelegramID, o.UserFirstName, username)

	if o.Payment.Status != domain.PaymentStatusPending {
		return core_tg_html.Sprintf(
			"💰 Оплата\n\n%s\n\nℹ️ Уже обработана: %s.",
			body,
			core_tg_html.Escape(paymentStatusText(o.Payment.Status)),
		)
	}

	return core_tg_html.Sprintf(
		"💰 Оплата на проверке\n\n%s\nОтправлена: %s",
		body,
		core_tg_html.Escape(o.Payment.UpdatedAt.In(time.Local).Format("02.01.2006 15:04")),
	)
}

func adminPaymentCardMarkup(p domain.Payment) tgbotapi.InlineKeyboardMarkup {
	backRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ К списку", "pay:admin"),
	)

	if p.Status != domain.PaymentStatusPending {
		return tgbotapi.NewInlineKeyboardMarkup(backRow)
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Подтвердить", paymentCallback("aconfirm", p.ID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отклонить", paymentCallback("areject", p.ID)),
		),
		backRow,
	)
}
