package access_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *AccessRepository) EnsureFeedToken(
	ctx context.Context,
	telegramID int64,
	newToken string,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO vpnbot.feed_tokens AS ft (user_id, token)
	SELECT u.id, $2
	FROM vpnbot.users u
	WHERE u.telegram_id = $1
	ON CONFLICT (user_id) DO UPDATE SET token = ft.token
	RETURNING ft.token;
	`

	var token string
	if err := r.pool.QueryRow(ctx, query, telegramID, newToken).Scan(&token); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return "", fmt.Errorf("user with telegram_id=%d: %w", telegramID, core_errors.ErrNotFound)
		}
		return "", fmt.Errorf("scan error: %w", err)
	}

	return token, nil
}
