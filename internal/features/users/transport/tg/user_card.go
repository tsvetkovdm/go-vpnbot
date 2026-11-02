package users_transport_tg

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
)

func assembleUserCard(
	o domain.UserOverview,
	filter domain.UserFilter,
	page int,
	now time.Time,
) (string, tgbotapi.InlineKeyboardMarkup) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("👤 %s\n", displayName(o.User)))
	sb.WriteString(fmt.Sprintf("🆔: %d\n\n", o.User.TelegramID))
	sb.WriteString(subscriptionLines(o, now))

	return sb.String(), userCardMarkup(o, filter, page)
}

func subscriptionLines(o domain.UserOverview, now time.Time) string {
	var sb strings.Builder
	sub := o.Subscription

	if sub != nil {
		if o.LastPlanTitle != nil {
			sb.WriteString(fmt.Sprintf("📦 Тариф: %s\n", *o.LastPlanTitle))
		}

		paidUntil := sub.PaidUntil(time.Local).Format("02.01.2006")
		if sub.HasAccess(now) {
			sb.WriteString(fmt.Sprintf("📅 Действует до: %s (%s)\n", paidUntil, daysLeftText(sub.DaysLeft(now))))
		} else {
			sb.WriteString(fmt.Sprintf("📅 Действовал до: %s\n", paidUntil))
		}
	}

	sb.WriteString(statusText(o, now))
	return sb.String()
}

func statusText(o domain.UserOverview, now time.Time) string {
	sub := o.Subscription

	switch {
	case o.User.Status == domain.StatusBlocked:
		return "🚫 Статус: заблокирован"
	case sub == nil:
		return "⚪ Статус: нет подписки"
	case !sub.HasAccess(now):
		return "🔴 Статус: Истёк"
	default:
		return "🟢 Статус: Активен"
	}
}

func daysLeftText(days int) string {
	if days == 0 {
		return "последний день"
	}
	return fmt.Sprintf("ещё %d дн.", days)
}

func userCardMarkup(o domain.UserOverview, filter domain.UserFilter, page int) tgbotapi.InlineKeyboardMarkup {
	userID := strconv.FormatInt(o.User.ID, 10)
	f := string(filter)
	p := strconv.Itoa(page)
	var rows [][]tgbotapi.InlineKeyboardButton

	if o.User.Status != domain.StatusBlocked {
		version := "0" // подписки ещё нет
		if o.Subscription != nil {
			version = strconv.FormatInt(o.Subscription.Version, 10)
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"➕ 30 дней",
				core_tg_callback.Build("subs", "extend", userID, "30", version, f, p),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"➕ 90 дней",
				core_tg_callback.Build("subs", "extend", userID, "90", version, f, p),
			),
		))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"✏️ Добавить дней",
				core_tg_callback.Build("subs", "days", userID, version, f, p),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"📅 Установить дату",
				core_tg_callback.Build("subs", "date", userID, version, f, p),
			),
		))
	}

	blockButton := tgbotapi.NewInlineKeyboardButtonData(
		"🚫 Заблокировать",
		core_tg_callback.Build("users", "block", userID, f, p),
	)
	if o.User.Status == domain.StatusBlocked {
		blockButton = tgbotapi.NewInlineKeyboardButtonData(
			"✅ Разблокировать",
			core_tg_callback.Build("users", "unblock", userID, f, p),
		)
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(blockButton))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", listCallback(filter, page)),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func parseUserCallback(data string) (int64, domain.UserFilter, int, error) {
	segments := core_tg_callback.Segments(data)
	if len(segments) != 5 {
		return 0, "", 0, fmt.Errorf("malformed callback data %q: %w", data, core_errors.ErrInvalidArgument)
	}

	userID, err := strconv.ParseInt(segments[2], 10, 64)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid userID %q: %w", segments[2], core_errors.ErrInvalidArgument)
	}

	filter := domain.UserFilter(segments[3])
	if _, ok := filterTitles[filter]; !ok {
		return 0, "", 0, fmt.Errorf("unknown filter %q: %w", segments[3], core_errors.ErrInvalidArgument)
	}

	page, err := strconv.Atoi(segments[4])
	if err != nil || page < 0 {
		return 0, "", 0, fmt.Errorf("invalid page %q: %w", segments[4], core_errors.ErrInvalidArgument)
	}

	return userID, filter, page, nil
}
