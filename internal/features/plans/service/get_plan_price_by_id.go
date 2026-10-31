package plans_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PlansService) GetPlanPriceByID(
	ctx context.Context,
	priceID int64,
) (domain.PlanPrice, error) {
	price, err := s.plansRepository.GetPlanPriceByID(ctx, priceID)
	if err != nil {
		return domain.PlanPrice{}, fmt.Errorf("get plan price: %w", err)
	}

	return price, nil
}
