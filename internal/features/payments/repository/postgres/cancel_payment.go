package payments_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *PaymentsRepository) CancelPayment(
	ctx context.Context,
	telegramID int64,
	paymentID int64,
) (domain.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE vpnbot.payments AS p
	SET status = 'cancelled', resolved_at = now(), updated_at = now()
	FROM vpnbot.users u
	WHERE p.id = $1
	  AND p.status = 'created'
	  AND u.id = p.user_id
	  AND u.telegram_id = $2
	RETURNING ` + paymentColumns + `;
	`

	m, err := scanPayment(r.pool.QueryRow(ctx, query, paymentID, telegramID))
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Payment{}, fmt.Errorf(
				"payment id=%d of telegram_id=%d is not cancellable: %w",
				paymentID,
				telegramID,
				core_errors.ErrConflict,
			)
		}
		return domain.Payment{}, fmt.Errorf("scan error: %w", err)
	}

	return m.toDomain(), nil
}
