package subscriptions_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *SubscriptionsRepository) UpdateSubscription(
	ctx context.Context,
	sub domain.Subscription,
) (domain.Subscription, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE vpnbot.subscriptions
	SET expires_at = $3, anchor_day = $4,
		notified_3_day_at = $5, notified_1_day_at = $6,
		updated_at = $7, version = version + 1
	WHERE id = $1 AND version = $2
	RETURNING ` + subscriptionColumns + `;`

	row := r.pool.QueryRow(
		ctx,
		query,
		sub.ID,
		sub.Version,
		sub.ExpiresAt,
		sub.AnchorDay,
		sub.Notified3DayAt,
		sub.Notified1DayAt,
		sub.UpdatedAt,
	)

	m, err := scanSubscription(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Subscription{}, fmt.Errorf(
				"subscription id=%d version=%d changed concurrently: %w",
				sub.ID, sub.Version, core_errors.ErrConflict,
			)
		}

		return domain.Subscription{}, fmt.Errorf("scan error: %w", err)
	}

	return m.toDomain(), nil
}
