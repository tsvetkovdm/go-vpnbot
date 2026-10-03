package plans_transport_tg

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
)

func parseCallbackData(data string) (int64, string, string, error) {
	segments := core_tg_callback.Segments(data)
	if len(segments) != 5 {
		return 0, "", "", fmt.Errorf("malformed callback data: %q", data)
	}

	planID, err := strconv.ParseInt(segments[2], 10, 64)
	if err != nil {
		return 0, "", "", fmt.Errorf("parse plan ID %q: %w", segments[2], err)
	}

	if segments[3] != "active" && segments[3] != "all" {
		return 0, "", "", fmt.Errorf("unknown filter: %q", segments[3])
	}

	return planID, segments[3], segments[4], nil
}

func assemblePlanCard(
	plan domain.Plan,
	filter string,
	page string,
) (string, tgbotapi.InlineKeyboardMarkup) {
	text := formatPlanInfo(plan)

	action, label := "activate", "Активировать"
	if plan.IsActive {
		action, label = "deactivate", "Деактивировать"
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				label,
				core_tg_callback.Build("plans", action, strconv.FormatInt(plan.ID, 10), filter, page),
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Назад",
				core_tg_callback.Build("plans", "list", filter, page),
			),
		),
	)

	return text, markup
}
