package subscriptions_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *SubscriptionsRepository) CreateSubscription(
	ctx context.Context,
	sub domain.Subscription,
) (domain.Subscription, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO vpnbot.subscriptions (user_id, expires_at, anchor_day, notified_3_day_at, notified_1_day_at, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING ` + subscriptionColumns + `;`

	row := r.pool.QueryRow(
		ctx,
		query,
		sub.UserID,
		sub.ExpiresAt,
		sub.AnchorDay,
		sub.Notified3DayAt,
		sub.Notified1DayAt,
		sub.CreatedAt,
		sub.UpdatedAt,
	)

	m, err := scanSubscription(row)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrUniqueViolation) {
			return domain.Subscription{}, fmt.Errorf(
				"subscription with user_id=%d unique violation: %w",
				sub.UserID,
				core_errors.ErrConflict,
			)
		}
		return domain.Subscription{}, fmt.Errorf("scan error: %w", err)
	}

	return m.toDomain(), nil
}
