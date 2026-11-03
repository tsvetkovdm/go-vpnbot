package access_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *AccessRepository) GetFeedOwnerByToken(
	ctx context.Context,
	token string,
) (domain.FeedOwner, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT u.id, u.xray_uuid, s.expires_at
	FROM vpnbot.feed_tokens ft
	JOIN vpnbot.users u ON u.id = ft.user_id
	LEFT JOIN vpnbot.subscriptions s ON s.user_id = u.id
	WHERE ft.token = $1;
	`

	var owner domain.FeedOwner
	err := r.pool.QueryRow(ctx, query, token).Scan(
		&owner.UserID,
		&owner.XrayUUID,
		&owner.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.FeedOwner{}, fmt.Errorf("feed token: %w", core_errors.ErrNotFound)
		}
		return domain.FeedOwner{}, fmt.Errorf("scan error: %w", err)
	}

	return owner, nil
}
