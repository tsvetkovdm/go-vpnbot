package subscriptions_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
)

const SkipVersionCheck int64 = -1

type SubscriptionsService struct {
	subscriptionsRepository SubscriptionsRepository
}

type SubscriptionsRepository interface {
	GetSubscriptionByUserID(
		ctx context.Context,
		userID int64,
	) (domain.Subscription, error)

	CreateSubscription(
		ctx context.Context,
		sub domain.Subscription,
	) (domain.Subscription, error)

	UpdateSubscription(
		ctx context.Context,
		sub domain.Subscription,
	) (domain.Subscription, error)

	GetUserIDByTelegramID(
		ctx context.Context,
		telegramID int64,
	) (int64, error)
}

func NewSubscriptionsService(
	subscriptionsRepository SubscriptionsRepository,
) *SubscriptionsService {
	return &SubscriptionsService{
		subscriptionsRepository: subscriptionsRepository,
	}
}

func (s *SubscriptionsService) loadOrNew(
	ctx context.Context,
	userID int64,
	now time.Time,
) (domain.Subscription, bool, error) {
	sub, err := s.subscriptionsRepository.GetSubscriptionByUserID(ctx, userID)
	if errors.Is(err, core_errors.ErrNotFound) {
		return domain.NewSubscription(userID, now), false, nil
	}
	if err != nil {
		return domain.Subscription{}, false, fmt.Errorf("get subscription by user_id=%d: %w", userID, err)
	}
	return sub, true, nil
}

func checkVersion(sub domain.Subscription, found bool, expected int64) error {
	current := int64(0)

	if found {
		current = sub.Version
	}

	if expected != SkipVersionCheck && expected != current {
		return fmt.Errorf("subscription version: expected %d, current %d: %w",
			expected, current, core_errors.ErrConflict)
	}

	return nil
}

func (s *SubscriptionsService) save(
	ctx context.Context,
	sub domain.Subscription,
	found bool,
) (domain.Subscription, error) {
	if found {
		return s.subscriptionsRepository.UpdateSubscription(ctx, sub)
	}

	return s.subscriptionsRepository.CreateSubscription(ctx, sub)
}
