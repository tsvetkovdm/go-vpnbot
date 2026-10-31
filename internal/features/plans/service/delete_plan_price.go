package plans_service

import (
	"context"
	"fmt"
)

func (s *PlansService) DeletePlanPrice(
	ctx context.Context,
	priceID int64,
) error {
	if err := s.plansRepository.DeletePlanPrice(ctx, priceID); err != nil {
		return fmt.Errorf("delete plan price: %w", err)
	}

	return nil
}
