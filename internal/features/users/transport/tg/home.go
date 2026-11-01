package users_transport_tg

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func assembleHome(o domain.UserOverview, now time.Time) (string, tgbotapi.InlineKeyboardMarkup) {
	plansRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("💳 Тарифы", "plans:public"),
	)

	if o.User.Status == domain.StatusBlocked {
		return "🚫 Доступ ограничен.\nЕсли это ошибка — напишите администратору.",
			homeMarkup(o, refreshRow())
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("👋 Привет, %s!\n\n", o.User.FirstName))

	sub := o.Subscription
	switch {
	case sub == nil:
		sb.WriteString("У вас пока нет доступа к VPN.\nВыберите тариф, чтобы подключиться.")
	case !sub.HasAccess(now):
		sb.WriteString(subscriptionLines(o, now))
		sb.WriteString("\n\nЧтобы продлить доступ, выберите тариф.")
	default:
		sb.WriteString(subscriptionLines(o, now))
	}

	rows := [][]tgbotapi.InlineKeyboardButton{}
	if row, ok := paymentRow(o); ok {
		rows = append(rows, row)
	}
	rows = append(rows, plansRow, refreshRow())

	return sb.String(), homeMarkup(o, rows...)
}

func paymentRow(o domain.UserOverview) ([]tgbotapi.InlineKeyboardButton, bool) {
	if o.UnresolvedPaymentStatus == nil {
		return nil, false
	}

	label := "💳 Продолжить оплату"
	if *o.UnresolvedPaymentStatus == domain.PaymentStatusPending {
		label = "⏳ Перевод на проверке"
	}

	return tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(label, "pay:current"),
	), true
}

func homeMarkup(
	o domain.UserOverview,
	rows ...[]tgbotapi.InlineKeyboardButton,
) tgbotapi.InlineKeyboardMarkup {
	if o.IsAdmin {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🛠 Админ-меню", "admins:menu"),
		))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func refreshRow() []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить", "users:home"),
	)
}
