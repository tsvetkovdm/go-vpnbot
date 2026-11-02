package plans_service

import (
	"context"
	"fmt"
)

func (s *PlansService) CanUseFreePrices(
	ctx context.Context,
	telegramID int64,
) (bool, error) {
	canUse, err := s.plansRepository.CanUseFreePrices(ctx, telegramID)
	if err != nil {
		return false, fmt.Errorf("can use free prices: %w", err)
	}

	return canUse, nil
}
