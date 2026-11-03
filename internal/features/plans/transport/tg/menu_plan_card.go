package plans_transport_tg

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
)

func loadPlanCard(
	ctx context.Context,
	plansService PlansService,
	planID int64,
	filter string,
	page string,
) (string, tgbotapi.InlineKeyboardMarkup, error) {
	plan, err := plansService.GetPlanByID(ctx, planID)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("get plan: %w", err)
	}

	prices, err := plansService.ListPlanPrices(ctx, planID)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("list plan prices: %w", err)
	}

	text, markup := assemblePlanCard(plan, prices, filter, page)
	return text, markup, nil
}

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

func backToListButton(filter, page string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(
		"⬅️ Назад",
		core_tg_callback.Build("plans", "list", filter, page),
	)
}

func viewPlanButton(planID int64, filter, page, label string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(
		label,
		core_tg_callback.Build("plans", "view", strconv.FormatInt(planID, 10), filter, page),
	)
}

func assemblePlanCard(
	plan domain.Plan,
	prices []domain.PlanPrice,
	filter string,
	page string,
) (string, tgbotapi.InlineKeyboardMarkup) {
	return formatPlanCard(plan, prices), planCardMarkup(plan, prices, filter, page)
}

func formatPlanCard(plan domain.Plan, prices []domain.PlanPrice) string {
	status := "🔴 Выключен"
	if plan.IsActive {
		status = "🟢 Включён"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "📦 Тариф «%s»\nСтатус: %s\n\n", plan.Title, status)

	if len(prices) == 0 {
		sb.WriteString("Сроков пока нет.\n")
	} else {
		sb.WriteString("Сроки:\n")
		for _, p := range prices {
			sb.WriteString(formatPriceLine(p))
			sb.WriteString("\n")
		}
	}

	if !hasActivePrice(prices) {
		sb.WriteString("\n⚠️ Нет активных сроков — пользователи тариф не увидят.")
	}

	return sb.String()
}

func hasActivePrice(prices []domain.PlanPrice) bool {
	for _, p := range prices {
		if p.IsActive {
			return true
		}
	}
	return false
}

func planCardMarkup(
	plan domain.Plan,
	prices []domain.PlanPrice,
	filter string,
	page string,
) tgbotapi.InlineKeyboardMarkup {
	id := strconv.FormatInt(plan.ID, 10)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(prices)+3)

	for _, p := range prices {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				formatPriceLine(p),
				core_tg_callback.Build("plans", "pview", strconv.FormatInt(p.ID, 10), filter, page),
			),
		))
	}

	toggleAction, toggleLabel := "activate", "✅ Включить"
	if plan.IsActive {
		toggleAction, toggleLabel = "deactivate", "⏸ Выключить"
	}

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"➕ Добавить срок",
				core_tg_callback.Build("plans", "addprice", id, filter, page),
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				toggleLabel,
				core_tg_callback.Build("plans", toggleAction, id, filter, page),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"🗑 Удалить",
				core_tg_callback.Build("plans", "delete", id, filter, page),
			),
		),
		tgbotapi.NewInlineKeyboardRow(backToListButton(filter, page)),
	)

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func plansMenuMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ В меню тарифов", core_tg_callback.Build("plans", "menu")),
		),
	)
}
