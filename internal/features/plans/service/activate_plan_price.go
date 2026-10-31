package plans_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PlansService) ActivatePlanPrice(
	ctx context.Context,
	priceID int64,
) (domain.PlanPrice, error) {
	planPrice, err := s.plansRepository.ActivatePlanPrice(ctx, priceID)
	if err != nil {
		return domain.PlanPrice{}, fmt.Errorf("activate plan price: %w", err)
	}

	return planPrice, nil
}
