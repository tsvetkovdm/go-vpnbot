package plans_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PlansService) ListPlanPrices(
	ctx context.Context,
	planID int64,
) ([]domain.PlanPrice, error) {
	planPrices, err := s.plansRepository.ListPlanPrices(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("get list of plan prices: %w", err)
	}

	return planPrices, nil
}
