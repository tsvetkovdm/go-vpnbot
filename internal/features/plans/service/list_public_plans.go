package plans_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PlansService) ListPublicPlans(
	ctx context.Context,
) ([]domain.Plan, error) {
	plans, err := s.plansRepository.ListPublicPlans(ctx)
	if err != nil {
		return nil, fmt.Errorf("list public plans: %w", err)
	}

	return plans, nil
}
