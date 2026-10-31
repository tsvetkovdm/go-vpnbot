package plans_transport_tg

import (
	"fmt"

	decimal "github.com/shopspring/decimal"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_text "github.com/tsvetkovdm/go-vpnbot/internal/core/text"
)

func formatPlanLine(plan domain.Plan) string {
	return fmt.Sprintf("#%d — %s %s", plan.ID, activeIcon(plan.IsActive), plan.Title)
}

func formatPriceForUser(price decimal.Decimal) string {
	if price.IsZero() {
		return "бесплатно"
	}
	if price.IsInteger() {
		return price.String() + " ₽"
	}
	return price.StringFixed(2) + " ₽"
}

func activeIcon(isActive bool) string {
	if isActive {
		return "🟢"
	}
	return "🔴"
}

func formatPriceDuration(d domain.Duration) string {
	if d.Months > 0 {
		n := int(d.Months)
		return fmt.Sprintf("%d %s", n, core_text.PluralRu(n, "месяц", "месяца", "месяцев"))
	}

	n := int(d.Days)
	return fmt.Sprintf("%d %s", n, core_text.PluralRu(n, "день", "дня", "дней"))
}

func formatPriceAmount(p domain.PlanPrice) string {
	amount := formatPriceForUser(p.Price)

	if p.Duration.Months < 2 || p.Price.IsZero() {
		return amount
	}

	months := decimal.NewFromInt(int64(p.Duration.Months))
	perMonth := p.Price.Div(months).Round(2)

	return fmt.Sprintf("%s (%s/мес)", amount, formatPriceForUser(perMonth))
}

func formatPriceLabel(p domain.PlanPrice) string {
	return formatPriceDuration(p.Duration) + " — " + formatPriceAmount(p)
}

func formatPriceLine(p domain.PlanPrice) string {
	return activeIcon(p.IsActive) + " " + formatPriceLabel(p)
}
