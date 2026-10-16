package subscriptions_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *SubscriptionsRepository) GetSubscriptionByUserID(
	ctx context.Context,
	userID int64,
) (domain.Subscription, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT ` + subscriptionColumns + ` 
	FROM vpnbot.subscriptions
	WHERE user_id = $1;
	`

	row := r.pool.QueryRow(ctx, query, userID)

	m, err := scanSubscription(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Subscription{}, fmt.Errorf(
				"subscription with user_id=%d: %w",
				userID,
				core_errors.ErrNotFound,
			)
		}

		return domain.Subscription{}, fmt.Errorf("scan error: %w", err)
	}

	return m.toDomain(), nil
}
