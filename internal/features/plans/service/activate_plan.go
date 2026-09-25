package plans_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PlansService) ActivatePlan(
	ctx context.Context,
	planID int64,
) (domain.Plan, error) {
	plan, err := s.plansRepository.ActivatePlan(ctx, planID)
	if err != nil {
		return domain.Plan{}, fmt.Errorf("activate plan: %w", err)
	}

	return plan, nil
}
