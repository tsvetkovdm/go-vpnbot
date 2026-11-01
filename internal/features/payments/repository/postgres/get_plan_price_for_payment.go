package payments_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *PaymentsRepository) GetPlanPriceForPayment(
	ctx context.Context,
	planPriceID int64,
) (domain.Plan, domain.PlanPrice, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT pp.id, pp.duration, pp.price, pp.is_active,
		   pl.id, pl.title, pl.is_active
	FROM vpnbot.plan_prices pp
	JOIN vpnbot.plans pl ON pl.id = pp.plan_id
	WHERE pp.id = $1;
	`

	var (
		price    domain.PlanPrice
		plan     domain.Plan
		duration pgtype.Interval
	)
	err := r.pool.QueryRow(ctx, query, planPriceID).Scan(
		&price.ID,
		&duration,
		&price.Price,
		&price.IsActive,
		&plan.ID,
		&plan.Title,
		&plan.IsActive,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Plan{}, domain.PlanPrice{}, fmt.Errorf(
				"plan price with id=%d: %w",
				planPriceID,
				core_errors.ErrNotFound,
			)
		}
		return domain.Plan{}, domain.PlanPrice{}, fmt.Errorf("scan error: %w", err)
	}

	price.PlanID = plan.ID
	price.Duration = durationFromInterval(duration)

	return plan, price, nil
}
