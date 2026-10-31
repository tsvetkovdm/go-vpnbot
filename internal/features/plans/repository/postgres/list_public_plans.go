package plans_postgres_repository

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (r *PlansRepository) ListPublicPlans(
	ctx context.Context,
) ([]domain.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT ` + planColumns + `
	FROM vpnbot.plans p
	WHERE p.is_active = true
	  AND EXISTS (
	      SELECT 1
	      FROM vpnbot.plan_prices pp
	      WHERE pp.plan_id = p.id AND pp.is_active = true
	  )
	ORDER BY p.id ASC;
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var plans []domain.Plan
	for rows.Next() {
		m, err := scanPlan(rows)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		plans = append(plans, m.toDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return plans, nil
}
