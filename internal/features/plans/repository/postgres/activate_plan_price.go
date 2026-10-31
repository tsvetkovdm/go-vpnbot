package plans_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *PlansRepository) ActivatePlanPrice(
	ctx context.Context,
	priceID int64,
) (domain.PlanPrice, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE vpnbot.plan_prices
	SET is_active = true, updated_at = now(), version = version + 1
	WHERE id = $1 AND is_active = false
	RETURNING ` + planPriceColumns + `;
	`

	row := r.pool.QueryRow(ctx, query, priceID)

	m, err := scanPlanPrice(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			if _, getErr := r.GetPlanPriceByID(ctx, priceID); getErr != nil {
				return domain.PlanPrice{}, getErr
			}
			return domain.PlanPrice{}, fmt.Errorf(
				"plan price id=%d already active: %w",
				priceID,
				core_errors.ErrConflict,
			)
		}
		if errors.Is(err, core_postgres_pool.ErrUniqueViolation) {
			return domain.PlanPrice{}, fmt.Errorf(
				"plan price id=%d: %w",
				priceID,
				domain.ErrPlanPriceDuplicateDuration,
			)
		}
		return domain.PlanPrice{}, fmt.Errorf("scan error: %w", err)
	}

	return m.toDomain(), nil
}
