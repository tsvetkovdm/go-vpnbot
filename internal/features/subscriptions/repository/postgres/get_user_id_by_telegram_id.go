package subscriptions_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *SubscriptionsRepository) GetUserIDByTelegramID(
	ctx context.Context,
	telegramID int64,
) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id 
	FROM vpnbot.users
	WHERE telegram_id = $1;
	`

	row := r.pool.QueryRow(ctx, query, telegramID)

	var userID int64
	if err := row.Scan(&userID); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return 0, fmt.Errorf(
				"user with telegram_id=%d: %w",
				telegramID,
				core_errors.ErrNotFound,
			)
		}

		return 0, fmt.Errorf("scan error: %w", err)
	}

	return userID, nil
}
