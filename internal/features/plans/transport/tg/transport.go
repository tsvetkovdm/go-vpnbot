package plans_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
	core_tg_middleware "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/middleware"
	core_tg_poller "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/poller"
)

type PlansTGHandler struct {
	bot            *tgbotapi.BotAPI
	plansService   PlansService
	adminOnly      core_tg_middleware.Middleware
	dialogStore    *core_tg_dialog.Store
	createPlanFlow core_tg_dialog.Flow
}

type PlansService interface {
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

	ActivatePlan(
		ctx context.Context,
		planID int64,
	) (domain.Plan, error)

	DeactivatePlan(
		ctx context.Context,
		planID int64,
	) (domain.Plan, error)
}

func NewPlansTGHandler(
	bot *tgbotapi.BotAPI,
	plansService PlansService,
	adminOnly core_tg_middleware.Middleware,
	dialogStore *core_tg_dialog.Store,
) *PlansTGHandler {
	return &PlansTGHandler{
		bot:            bot,
		plansService:   plansService,
		adminOnly:      adminOnly,
		dialogStore:    dialogStore,
		createPlanFlow: newCreatePlanFlow(plansService),
	}
}

func (h *PlansTGHandler) Routes() []core_tg_poller.Route {
	return []core_tg_poller.Route{
		core_tg_poller.NewRoute(
			"create_plan",
			core_tg_poller.TriggerCommand,
			h.CreatePlan,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"get_plan",
			core_tg_poller.TriggerCommand,
			h.GetPlan,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"list_plans",
			core_tg_poller.TriggerCommand,
			h.ListPlans,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"deactivate_plan",
			core_tg_poller.TriggerCommand,
			h.DeactivatePlan,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"activate_plan",
			core_tg_poller.TriggerCommand,
			h.ActivatePlan,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"plans:menu",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuRoot,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"plans:create",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuCreatePlan,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"plans:list:menu",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuListChoice,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"plans:list:active:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuListPlans,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"plans:list:all:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuListPlans,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"plans:view:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuViewPlan,
			h.adminOnly,
		),
	}
}

func (h *PlansTGHandler) Flows() []core_tg_dialog.Flow {
	return []core_tg_dialog.Flow{h.createPlanFlow}
}
