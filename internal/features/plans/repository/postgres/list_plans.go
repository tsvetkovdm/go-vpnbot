package plans_postgres_repository

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (r *PlansRepository) ListPlans(
	ctx context.Context,
	includeInactive bool,
) ([]domain.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	where := "is_active = true"
	if includeInactive {
		where = "TRUE"
	}

	query := `
	SELECT ` + planColumns + `
	FROM vpnbot.plans
	WHERE ` + where + `
	ORDER BY id ASC;
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
