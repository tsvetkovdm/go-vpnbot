package payments_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *PaymentsRepository) GetUserStatusByTelegramID(
	ctx context.Context,
	telegramID int64,
) (int64, domain.Status, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, status
	FROM vpnbot.users
	WHERE telegram_id = $1;
	`

	var (
		userID int64
		status string
	)
	err := r.pool.QueryRow(ctx, query, telegramID).Scan(&userID, &status)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return 0, "", fmt.Errorf(
				"user with telegram_id=%d: %w",
				telegramID,
				core_errors.ErrNotFound,
			)
		}
		return 0, "", fmt.Errorf("scan error: %w", err)
	}

	return userID, domain.Status(status), nil
}
