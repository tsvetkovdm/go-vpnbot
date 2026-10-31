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

	CreatePlanPrice(
		ctx context.Context,
		price domain.PlanPrice,
	) (domain.PlanPrice, error)

	GetPlanByID(
		ctx context.Context,
		planID int64,
	) (domain.Plan, error)

	GetPlanPriceByID(
		ctx context.Context,
		priceID int64,
	) (domain.PlanPrice, error)

	ListPlans(
		ctx context.Context,
		includeInactive bool,
	) ([]domain.Plan, error)

	ListPlanPrices(
		ctx context.Context,
		planID int64,
	) ([]domain.PlanPrice, error)

	ActivatePlan(
		ctx context.Context,
		planID int64,
	) (domain.Plan, error)

	DeactivatePlan(
		ctx context.Context,
		planID int64,
	) (domain.Plan, error)

	DeletePlan(
		ctx context.Context,
		planID int64,
	) error
}

func NewPlansService(
	plansRepository PlansRepository,
) *PlansService {
	return &PlansService{
		plansRepository: plansRepository,
	}
}
