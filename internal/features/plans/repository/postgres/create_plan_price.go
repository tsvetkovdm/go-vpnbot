package plans_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *PlansRepository) CreatePlanPrice(
	ctx context.Context,
	price domain.PlanPrice,
) (domain.PlanPrice, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO vpnbot.plan_prices (plan_id, duration, price, is_active)
	VALUES ($1, $2, $3, $4)
	RETURNING ` + planPriceColumns + `;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		price.PlanID,
		intervalFromDuration(price.Duration),
		price.Price,
		price.IsActive,
	)

	m, err := scanPlanPrice(row)
	if err != nil {
		switch {
		case errors.Is(err, core_postgres_pool.ErrUniqueViolation):
			return domain.PlanPrice{}, fmt.Errorf(
				"active price with duration %+v already exists in plan id=%d: %w",
				price.Duration,
				price.PlanID,
				core_errors.ErrConflict,
			)
		case errors.Is(err, core_postgres_pool.ErrViolatesForeignKey):
			return domain.PlanPrice{}, fmt.Errorf(
				"plan with id=%d: %w",
				price.PlanID,
				core_errors.ErrNotFound,
			)
		}
		return domain.PlanPrice{}, fmt.Errorf("scan error: %w", err)
	}

	return m.toDomain(), nil
}
