package plans_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PlansService) CreatePlan(
	ctx context.Context,
	plan domain.Plan,
) (domain.Plan, error) {
	plan, err := s.plansRepository.CreatePlan(ctx, plan)
	if err != nil {
		return domain.Plan{}, fmt.Errorf("create plan: %w", err)
	}

	return plan, nil
}
