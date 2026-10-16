package users_transport_tg

import (
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
)

const usersPageSize = 8

var filterTitles = map[domain.UserFilter]string{
	domain.UserFilterActive:   "✅ С доступом",
	domain.UserFilterExpiring: "⏳ Истекает ≤ 7 дн",
	domain.UserFilterInactive: "❌ Без доступа",
	domain.UserFilterBlocked:  "🚫 Заблокированы",
	domain.UserFilterAll:      "👥 Все",
}

var filterOrder = []domain.UserFilter{
	domain.UserFilterActive,
	domain.UserFilterExpiring,
	domain.UserFilterInactive,
	domain.UserFilterBlocked,
	domain.UserFilterAll,
}

func displayName(u domain.User) string {
	if u.Username != nil {
		return fmt.Sprintf("%s (@%s)", u.FirstName, *u.Username)
	}

	return u.FirstName
}

func userButtonLabel(o domain.UserOverview, now time.Time) string {
	name := displayName(o.User)

	if o.User.Status == domain.StatusBlocked {
		return "🚫 " + name
	}

	if o.Subscription != nil && o.Subscription.HasAccess(now) {
		paidUntil := o.Subscription.PaidUntil(time.Local).Format("02.01")
		return fmt.Sprintf("✅ %s · до %s", name, paidUntil)
	}

	return fmt.Sprintf("❌ %s · нет доступа", name)
}

func listCallback(filter domain.UserFilter, page int) string {
	return core_tg_callback.Build("users", "list", string(filter), strconv.Itoa(page))
}

func backToUsersMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "users:menu"),
		),
	)
}
