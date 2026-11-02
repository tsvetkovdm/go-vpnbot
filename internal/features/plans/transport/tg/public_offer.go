package plans_transport_tg

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
)

var errPlanUnavailable = fmt.Errorf("plan is unavailable: %w", core_errors.ErrNotFound)

func parseIDArg(data string) (int64, error) {
	segments := core_tg_callback.Segments(data)
	if len(segments) != 3 {
		return 0, fmt.Errorf("malformed callback data: %q", data)
	}

	id, err := strconv.ParseInt(segments[2], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse ID %q: %w", segments[2], err)
	}

	return id, nil
}

func publicPrices(prices []domain.PlanPrice, includeFree bool) []domain.PlanPrice {
	var result []domain.PlanPrice
	for _, p := range prices {
		if !p.IsActive {
			continue
		}
		if p.Price.IsZero() && !includeFree {
			continue
		}
		result = append(result, p)
	}
	return result
}

func loadPublicOffer(
	ctx context.Context,
	plansService PlansService,
	plan domain.Plan,
	backData string,
	includeFree bool,
) (string, tgbotapi.InlineKeyboardMarkup, error) {
	if !plan.IsActive {
		return "", tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("plan id=%d: %w", plan.ID, errPlanUnavailable)
	}

	prices, err := plansService.ListPlanPrices(ctx, plan.ID)
	if err != nil {
		return "", tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("list plan prices: %w", err)
	}

	prices = publicPrices(prices, includeFree)
	if len(prices) == 0 {
		return "", tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("plan id=%d has no active prices: %w", plan.ID, errPlanUnavailable)
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(prices)+1)
	for _, p := range prices {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				formatPriceLabel(p),
				core_tg_callback.Build("plans", "buy", strconv.FormatInt(p.ID, 10)),
			),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", backData),
	))

	text := fmt.Sprintf("💳 Тариф «%s»\n\nВыберите срок:", plan.Title)
	return text, tgbotapi.NewInlineKeyboardMarkup(rows...), nil
}
