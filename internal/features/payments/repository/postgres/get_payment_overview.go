package payments_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *PaymentsRepository) GetPaymentOverview(
	ctx context.Context,
	paymentID int64,
) (domain.PaymentOverview, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT ` + paymentOverviewColumns + `
	FROM vpnbot.payments p
	JOIN vpnbot.users u ON u.id = p.user_id
	WHERE p.id = $1;
	`

	overview, err := scanPaymentOverview(r.pool.QueryRow(ctx, query, paymentID))
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.PaymentOverview{}, fmt.Errorf(
				"payment with id=%d: %w",
				paymentID,
				core_errors.ErrNotFound,
			)
		}
		return domain.PaymentOverview{}, fmt.Errorf("scan error: %w", err)
	}

	return overview, nil
}
