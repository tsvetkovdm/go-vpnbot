package plans_postgres_repository

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (r *PlansRepository) DeactivatePlan(
	ctx context.Context,
	plan domain.Plan,
) (domain.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE vpnbot.plans
	SET is_active = false
	WHERE id = $1 AND version = $2
	RETURNING id, version, title, duration, price, is_active, created_at, updated_at;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		plan.ID,
		plan.Version,
	)

	var planModel PlanModel
	err := row.Scan(
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
		return domain.Plan{}, fmt.Errorf("scan error: %w", err)
	}

	return planModel.toDomain(), nil
}
