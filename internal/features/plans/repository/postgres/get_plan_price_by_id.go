package plans_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *PlansRepository) GetPlanPriceByID(
	ctx context.Context,
	priceID int64,
) (domain.PlanPrice, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT ` + planPriceColumns + `
	FROM vpnbot.plan_prices
	WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, priceID)

	m, err := scanPlanPrice(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.PlanPrice{}, fmt.Errorf(
				"plan price with id=%d: %w",
				priceID,
				core_errors.ErrNotFound,
			)
		}

		return domain.PlanPrice{}, fmt.Errorf("scan error: %w", err)
	}

	return m.toDomain(), nil
}
