package plans_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *PlansRepository) GetPlanByID(
	ctx context.Context,
	planID int64,
) (domain.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT ` + planColumns + `
	FROM vpnbot.plans
	WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, planID)

	m, err := scanPlan(row)
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

	return m.toDomain(), nil
}
