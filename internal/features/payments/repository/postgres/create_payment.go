package payments_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *PaymentsRepository) CreatePayment(
	ctx context.Context,
	payment domain.Payment,
) (domain.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO vpnbot.payments AS p (
		user_id, plan_price_id,
		plan_title_snapshot, price_snapshot, duration_snapshot,
		comment_code, status
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING ` + paymentColumns + `;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		payment.UserID,
		payment.PlanPriceID,
		payment.PlanTitle,
		payment.Price,
		intervalFromDuration(payment.Duration),
		payment.CommentCode,
		string(payment.Status),
	)

	m, err := scanPayment(row)
	if err != nil {
		switch {
		case errors.Is(err, core_postgres_pool.ErrUniqueViolation):
			return domain.Payment{}, fmt.Errorf(
				"user id=%d: %w",
				payment.UserID,
				domain.ErrPaymentUnresolvedExists,
			)
		case errors.Is(err, core_postgres_pool.ErrViolatesForeignKey):
			return domain.Payment{}, fmt.Errorf(
				"user id=%d or plan price id=%d: %w",
				payment.UserID,
				payment.PlanPriceID,
				core_errors.ErrNotFound,
			)
		}
		return domain.Payment{}, fmt.Errorf("scan error: %w", err)
	}

	return m.toDomain(), nil
}
