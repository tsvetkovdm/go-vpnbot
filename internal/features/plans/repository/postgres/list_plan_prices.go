package plans_postgres_repository

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (r *PlansRepository) ListPlanPrices(
	ctx context.Context,
	planID int64,
) ([]domain.PlanPrice, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT ` + planPriceColumns + `
	FROM vpnbot.plan_prices
 	WHERE plan_id = $1
 	ORDER BY duration, id;
	`

	rows, err := r.pool.Query(ctx, query, planID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var prices []domain.PlanPrice
	for rows.Next() {
		m, err := scanPlanPrice(rows)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		prices = append(prices, m.toDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return prices, nil
}
