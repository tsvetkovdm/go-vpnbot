package plans_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PlansService) DeactivatePlan(
	ctx context.Context,
	planID int64,
) (domain.Plan, error) {
	plan, err := s.plansRepository.DeactivatePlan(ctx, planID)
	if err != nil {
		return domain.Plan{}, fmt.Errorf("deactivate plan: %w", err)
	}

	return plan, nil
}
