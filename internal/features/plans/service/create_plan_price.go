package plans_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PlansService) CreatePlanPrice(
	ctx context.Context,
	price domain.PlanPrice,
) (domain.PlanPrice, error) {
	price, err := s.plansRepository.CreatePlanPrice(ctx, price)
	if err != nil {
		return domain.PlanPrice{}, fmt.Errorf("create plan price: %w", err)
	}

	return price, nil
}
