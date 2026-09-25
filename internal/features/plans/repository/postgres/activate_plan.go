package plans_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *PlansRepository) ActivatePlan(
	ctx context.Context,
	planID int64,
) (domain.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE vpnbot.plans
	SET is_active = true, updated_at = $2, version = version + 1
	WHERE id = $1
	RETURNING id, version, title, duration, price, is_active, created_at, updated_at;
	`

	row := r.pool.QueryRow(ctx, query, planID, time.Now())

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
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Plan{}, fmt.Errorf(
				"plan with id=%d: %w",
				planID,
				core_errors.ErrNotFound,
			)
		}
		return domain.Plan{}, fmt.Errorf("scan error: %w", err)
	}

	return planModel.toDomain(), nil
}
