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
	sb.WriteString(fmt.Sprintf("Telegram ID: %d\n\n", o.User.TelegramID))
	sb.WriteString(accessStatus(o, now))

	return sb.String(), userCardMarkup(o, filter, page)
}

func accessStatus(o domain.UserOverview, now time.Time) string {
	sub := o.Subscription

	switch {
	case o.User.Status == domain.StatusBlocked:
		return "🚫 Заблокирован"
	case sub == nil:
		return "❌ Доступа не было"
	case !sub.HasAccess(now):
		return "❌ Доступ закончился " + sub.PaidUntil(time.Local).Format("02.01.2006")
	}

	left := "последний день"
	if days := sub.DaysLeft(now); days > 0 {
		left = fmt.Sprintf("осталось %d дн.", days)
	}

	return fmt.Sprintf(
		"✅ Доступ до %s (%s)\n📅 Дата оплаты: %d-е число",
		sub.PaidUntil(time.Local).Format("02.01.2006"),
		left,
		sub.AnchorDay,
	)
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
