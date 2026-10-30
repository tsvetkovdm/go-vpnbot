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
			tgbotapi.NewInlineKeyboardMarkup(refreshRow())
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

	return sb.String(), tgbotapi.NewInlineKeyboardMarkup(plansRow, refreshRow())
}

func refreshRow() []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить", "users:home"),
	)
}
