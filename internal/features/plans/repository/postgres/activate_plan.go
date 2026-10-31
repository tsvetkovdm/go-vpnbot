package plans_postgres_repository

import (
	"context"
	"errors"
	"fmt"

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
	SET is_active = true, updated_at = now(), version = version + 1
	WHERE id = $1 AND is_active = false
	RETURNING ` + planColumns + `;
	`

	row := r.pool.QueryRow(ctx, query, planID)

	m, err := scanPlan(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			if _, getErr := r.GetPlanByID(ctx, planID); getErr != nil {
				return domain.Plan{}, getErr
			}
			return domain.Plan{}, fmt.Errorf(
				"plan id=%d already active: %w",
				planID,
				core_errors.ErrConflict,
			)
		}
		return domain.Plan{}, fmt.Errorf("scan error: %w", err)
	}

	return m.toDomain(), nil
}
