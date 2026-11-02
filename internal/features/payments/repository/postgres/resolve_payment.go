package payments_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *PaymentsRepository) ResolvePayment(
	ctx context.Context,
	paymentID int64,
	adminTelegramID int64,
	status domain.PaymentStatus,
) (domain.Payment, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE vpnbot.payments AS p
	SET status = $3,
		resolved_by_admin_id = (
			SELECT a.id FROM vpnbot.admins a
			WHERE a.telegram_id = $2 AND a.removed_at IS NULL
		),
		resolved_at = now(),
		updated_at = now()
	FROM vpnbot.users u
	WHERE p.id = $1
	  AND p.status = 'pending'
	  AND u.id = p.user_id
	RETURNING ` + paymentColumns + `, u.telegram_id;
	`

	var userTelegramID int64
	m, err := scanPayment(
		r.pool.QueryRow(ctx, query, paymentID, adminTelegramID, string(status)),
		&userTelegramID,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Payment{}, 0, fmt.Errorf(
				"payment id=%d is not pending: %w",
				paymentID,
				domain.ErrPaymentAlreadyResolved,
			)
		}
		return domain.Payment{}, 0, fmt.Errorf("scan error: %w", err)
	}

	return m.toDomain(), userTelegramID, nil
}
