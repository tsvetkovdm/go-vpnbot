package subscriptions_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *SubscriptionsService) ExtendDays(
	ctx context.Context,
	userID int64,
	days int,
	expectedVersion int64,
) (domain.Subscription, error) {
	if err := domain.ValidateExtendDays(days); err != nil {
		return domain.Subscription{}, fmt.Errorf("invalid days value: %w", err)
	}

	sub, err := s.Extend(ctx, userID, domain.Duration{Days: int32(days)}, expectedVersion)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("extend days: %w", err)
	}

	return sub, nil
}
