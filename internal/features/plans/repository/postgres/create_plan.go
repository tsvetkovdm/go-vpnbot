package plans_postgres_repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (r *PlansRepository) CreatePlan(
	ctx context.Context,
	plan domain.Plan,
) (domain.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO vpnbot.plans (title, duration, price, is_active, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, version, title, duration, price, is_active, created_at, updated_at;
	`

	duration := pgtype.Interval{
		Months:       plan.Duration.Months,
		Days:         plan.Duration.Days,
		Microseconds: plan.Duration.Microseconds,
		Valid:        true,
	}

	row := r.pool.QueryRow(
		ctx,
		query,
		plan.Title,
		duration,
		plan.Price,
		plan.IsActive,
		plan.CreatedAt,
		plan.UpdatedAt,
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
