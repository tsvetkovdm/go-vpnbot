package plans_service

import (
	"context"
	"fmt"
)

func (s *PlansService) DeletePlan(
	ctx context.Context,
	planID int64,
) error {
	if err := s.plansRepository.DeletePlan(ctx, planID); err != nil {
		return fmt.Errorf("delete plan: %w", err)
	}

	return nil
}
