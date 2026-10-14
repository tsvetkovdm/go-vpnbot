package plans_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *PlansRepository) DeletePlan(
	ctx context.Context,
	planID int64,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM vpnbot.plans
	WHERE id = $1;
	`
	cmdTag, err := r.pool.Exec(ctx, query, planID)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return fmt.Errorf(
				"plan id=%d is in use: %w",
				planID,
				core_errors.ErrConflict,
			)
		}
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"plan with id='%d': %w",
			planID,
			core_errors.ErrNotFound,
		)
	}

	return nil
}
