package payments_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *PaymentsRepository) GetPaymentByID(
	ctx context.Context,
	paymentID int64,
) (domain.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT ` + paymentColumns + `
	FROM vpnbot.payments p
	WHERE p.id = $1;
	`

	m, err := scanPayment(r.pool.QueryRow(ctx, query, paymentID))
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Payment{}, fmt.Errorf(
				"payment with id=%d: %w",
				paymentID,
				core_errors.ErrNotFound,
			)
		}
		return domain.Payment{}, fmt.Errorf("scan error: %w", err)
	}

	return m.toDomain(), nil
}
