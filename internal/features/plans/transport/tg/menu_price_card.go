package plans_transport_tg

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
)

func loadPriceCard(
	ctx context.Context,
	plansService PlansService,
	priceID int64,
	filter string,
	page string,
) (string, tgbotapi.InlineKeyboardMarkup, error) {
	price, err := plansService.GetPlanPriceByID(ctx, priceID)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("get plan price: %w", err)
	}

	plan, err := plansService.GetPlanByID(ctx, price.PlanID)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("get plan: %w", err)
	}

	return formatPriceCard(plan, price), priceCardMarkup(price, filter, page), nil
}

func formatPriceCard(plan domain.Plan, price domain.PlanPrice) string {
	status := "🔴 Выключен"
	if price.IsActive {
		status = "🟢 Включён"
	}

	return fmt.Sprintf(
		"📦 Тариф «%s»\n\nСрок: %s\nЦена: %s\nСтатус: %s",
		plan.Title,
		formatPriceDuration(price.Duration),
		formatPriceAmount(price),
		status,
	)
}

func priceCardMarkup(price domain.PlanPrice, filter, page string) tgbotapi.InlineKeyboardMarkup {
	id := strconv.FormatInt(price.ID, 10)

	toggleAction, toggleLabel := "pactivate", "✅ Включить"
	if price.IsActive {
		toggleAction, toggleLabel = "pdeactivate", "⏸ Выключить"
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				toggleLabel,
				core_tg_callback.Build("plans", toggleAction, id, filter, page),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"🗑 Удалить",
				core_tg_callback.Build("plans", "pdelete", id, filter, page),
			),
		),
		tgbotapi.NewInlineKeyboardRow(viewPlanButton(price.PlanID, filter, page, "⬅️ К тарифу")),
	)
}
