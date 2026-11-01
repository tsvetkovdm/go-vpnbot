package payments_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *PaymentsRepository) GetUnresolvedPaymentByTelegramID(
	ctx context.Context,
	telegramID int64,
) (domain.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT ` + paymentColumns + `
	FROM vpnbot.payments p
	JOIN vpnbot.users u ON u.id = p.user_id
	WHERE u.telegram_id = $1
	  AND p.status IN ('created', 'pending');
	`

	m, err := scanPayment(r.pool.QueryRow(ctx, query, telegramID))
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Payment{}, fmt.Errorf(
				"unresolved payment for telegram_id=%d: %w",
				telegramID,
				core_errors.ErrNotFound,
			)
		}
		return domain.Payment{}, fmt.Errorf("scan error: %w", err)
	}

	return m.toDomain(), nil
}
