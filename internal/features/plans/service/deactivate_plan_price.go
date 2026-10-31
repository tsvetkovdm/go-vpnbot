package plans_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PlansService) DeactivatePlanPrice(
	ctx context.Context,
	priceID int64,
) (domain.PlanPrice, error) {
	planPrice, err := s.plansRepository.DeactivatePlanPrice(ctx, priceID)
	if err != nil {
		return domain.PlanPrice{}, fmt.Errorf("deactivate plan price: %w", err)
	}

	return planPrice, nil
}
