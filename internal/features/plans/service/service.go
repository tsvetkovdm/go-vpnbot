package plans_service

import (
	"context"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

type PlansService struct {
	plansRepository PlansRepository
}

type PlansRepository interface {
	CreatePlan(
		ctx context.Context,
		plan domain.Plan,
	) (domain.Plan, error)

	GetPlanByID(
		ctx context.Context,
		planID int64,
	) (domain.Plan, error)

	ListPlans(
		ctx context.Context,
		includeInactive bool,
	) ([]domain.Plan, error)
}

func NewPlansService(
	plansRepository PlansRepository,
) *PlansService {
	return &PlansService{
		plansRepository: plansRepository,
	}
}
