package subscriptions_service

import (
	"context"
	"fmt"
)

func (s *SubscriptionsService) ResolveUserID(
	ctx context.Context,
	telegramID int64,
) (int64, error) {
	id, err := s.subscriptionsRepository.GetUserIDByTelegramID(ctx, telegramID)
	if err != nil {
		return 0, fmt.Errorf("resolve user with telegramID=%d: %w", telegramID, err)
	}

	return id, nil
}
