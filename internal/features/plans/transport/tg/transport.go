package plans_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
	core_tg_middleware "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/middleware"
	core_tg_poller "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/poller"
	core_tg_screen "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/screen"
)

type PlansTGHandler struct {
	bot            *tgbotapi.BotAPI
	plansService   PlansService
	adminOnly      core_tg_middleware.Middleware
	dialogStore    *core_tg_dialog.Store
	createPlanFlow core_tg_dialog.Flow
	addPriceFlow   core_tg_dialog.Flow
	screens        *core_tg_screen.Screens
}

type PlansService interface {
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

	ListPublicPlans(
		ctx context.Context,
		includeFree bool,
	) ([]domain.Plan, error)

	CanUseFreePrices(
		ctx context.Context,
		telegramID int64,
	) (bool, error)

	ActivatePlan(
		ctx context.Context,
		planID int64,
	) (domain.Plan, error)

	ActivatePlanPrice(
		ctx context.Context,
		priceID int64,
	) (domain.PlanPrice, error)

	DeactivatePlan(
		ctx context.Context,
		planID int64,
	) (domain.Plan, error)

	DeactivatePlanPrice(
		ctx context.Context,
		priceID int64,
	) (domain.PlanPrice, error)

	DeletePlan(
		ctx context.Context,
		planID int64,
	) error

	DeletePlanPrice(
		ctx context.Context,
		priceID int64,
	) error
}

func NewPlansTGHandler(
	bot *tgbotapi.BotAPI,
	plansService PlansService,
	adminOnly core_tg_middleware.Middleware,
	dialogStore *core_tg_dialog.Store,
	screens *core_tg_screen.Screens,
) *PlansTGHandler {
	return &PlansTGHandler{
		bot:            bot,
		plansService:   plansService,
		adminOnly:      adminOnly,
		dialogStore:    dialogStore,
		createPlanFlow: newCreatePlanFlow(plansService),
		addPriceFlow:   newAddPriceFlow(plansService),
		screens:        screens,
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
			"delete_plan",
			core_tg_poller.TriggerCommand,
			h.DeletePlan,
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
			"plans:addprice:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuAddPrice,
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
		core_tg_poller.NewRoute(
			"plans:pview:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuViewPrice,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"plans:activate:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuTogglePlan,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"plans:deactivate:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuTogglePlan,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"plans:pactivate:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuTogglePrice,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"plans:pdeactivate:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuTogglePrice,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"plans:delete:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuDeletePlanAsk,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"plans:confirm_delete:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuDeletePlanConfirm,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"plans:pdelete:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuDeletePriceAsk,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"plans:pconfirm_delete:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuDeletePriceConfirm,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"plans:public",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuPublicPlans,
		),
		core_tg_poller.NewRoute(
			"plans:offer:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuPublicOffer,
		),
		core_tg_poller.NewRoute(
			"plans:buy:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuPublicBuy,
		),
	}
}

func (h *PlansTGHandler) Flows() []core_tg_dialog.Flow {
	return []core_tg_dialog.Flow{h.createPlanFlow, h.addPriceFlow}
}
