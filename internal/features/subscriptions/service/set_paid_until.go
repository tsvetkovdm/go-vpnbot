package subscriptions_service

import (
	"context"
	"fmt"
	"time"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *SubscriptionsService) SetPaidUntil(
	ctx context.Context,
	userID int64,
	date time.Time,
	expectedVersion int64,
) (domain.Subscription, error) {
	now := time.Now()

	sub, found, err := s.loadOrNew(ctx, userID, now)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("set paid until: %w", err)
	}

	if err := checkVersion(sub, found, expectedVersion); err != nil {
		return domain.Subscription{}, err
	}

	sub, err = sub.SetPaidUntil(date, now)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("set paid until: %w", err)
	}

	sub, err = s.save(ctx, sub, found)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("set paid until: %w", err)
	}

	return sub, nil
}
