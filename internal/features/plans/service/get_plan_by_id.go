package plans_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PlansService) GetPlanByID(
	ctx context.Context,
	planID int64,
) (domain.Plan, error) {
	plan, err := s.plansRepository.GetPlanByID(ctx, planID)
	if err != nil {
		return domain.Plan{}, fmt.Errorf("get plan from repository: %w", err)
	}

	return plan, nil
}
