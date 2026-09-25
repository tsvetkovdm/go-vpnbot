package plans_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PlansService) ListPlans(
	ctx context.Context,
	includeInactive bool,
) ([]domain.Plan, error) {
	plans, err := s.plansRepository.ListPlans(ctx, includeInactive)
	if err != nil {
		return nil, fmt.Errorf("get list of plans: %w", err)
	}

	return plans, nil
}
