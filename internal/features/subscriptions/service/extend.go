package subscriptions_service

import (
	"context"
	"fmt"
	"time"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *SubscriptionsService) Extend(
	ctx context.Context,
	userID int64,
	d domain.Duration,
	expectedVersion int64,
) (domain.Subscription, error) {
	now := time.Now()

	sub, found, err := s.loadOrNew(ctx, userID, now)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("extend subscription: %w", err)
	}

	if err := checkVersion(sub, found, expectedVersion); err != nil {
		return domain.Subscription{}, fmt.Errorf("extend subscription: %w", err)
	}

	sub = sub.Extend(d, now)

	sub, err = s.save(ctx, sub, found)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("extend subscription: %w", err)
	}

	return sub, nil
}
