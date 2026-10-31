package plans_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *PlansRepository) DeletePlanPrice(
	ctx context.Context,
	priceID int64,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM vpnbot.plan_prices
	WHERE id = $1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, priceID)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return fmt.Errorf(
				"plan price id=%d is in use: %w",
				priceID,
				core_errors.ErrConflict,
			)
		}
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"plan price with id=%d: %w",
			priceID,
			core_errors.ErrNotFound,
		)
	}

	return nil
}
