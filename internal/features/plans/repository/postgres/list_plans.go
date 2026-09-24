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

	query := `
	SELECT id, version, title, duration, price, is_active, created_at, updated_at
	FROM vpnbot.plans
	WHERE is_active = true
	ORDER BY id ASC;
	`
	if includeInactive {
		query = `
		SELECT id, version, title, duration, price, is_active, created_at, updated_at
		FROM vpnbot.plans
		ORDER BY id ASC;
		`
	}

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var plans []domain.Plan
	for rows.Next() {
		var planModel PlanModel

		err := rows.Scan(
			&planModel.ID,
			&planModel.Version,
			&planModel.Title,
			&planModel.Duration,
			&planModel.Price,
			&planModel.IsActive,
			&planModel.CreatedAt,
			&planModel.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		plans = append(plans, planModel.toDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return plans, nil
}
