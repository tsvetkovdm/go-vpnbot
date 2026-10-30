package subscriptions_transport_tg

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
	core_tg_middleware "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/middleware"
	core_tg_poller "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/poller"
)

type SubscriptionsTGHandler struct {
	bot                  *tgbotapi.BotAPI
	subscriptionsService SubscriptionsService
	adminOnly            core_tg_middleware.Middleware
	dialogStore          *core_tg_dialog.Store
	extendDaysFlow       core_tg_dialog.Flow
	setDateFlow          core_tg_dialog.Flow
}

type SubscriptionsService interface {
	ExtendDays(
		ctx context.Context,
		userID int64,
		days int,
		expectedVersion int64,
	) (domain.Subscription, error)

	SetPaidUntil(
		ctx context.Context,
		userID int64,
		date time.Time,
		expectedVersion int64,
	) (domain.Subscription, error)

	ResolveUserID(
		ctx context.Context,
		telegramID int64,
	) (int64, error)
}

func NewSubscriptionsTGHandler(
	bot *tgbotapi.BotAPI,
	subscriptionsService SubscriptionsService,
	adminOnly core_tg_middleware.Middleware,
	dialogStore *core_tg_dialog.Store,
) *SubscriptionsTGHandler {
	return &SubscriptionsTGHandler{
		bot:                  bot,
		subscriptionsService: subscriptionsService,
		adminOnly:            adminOnly,
		dialogStore:          dialogStore,
		extendDaysFlow:       newExtendDaysFlow(subscriptionsService),
		setDateFlow:          newSetDateFlow(subscriptionsService),
	}
}

func (h *SubscriptionsTGHandler) Routes() []core_tg_poller.Route {
	return []core_tg_poller.Route{
		core_tg_poller.NewRoute(
			"extend",
			core_tg_poller.TriggerCommand,
			h.ExtendCommand,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"subs:extend:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuExtendDays,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"subs:days:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuStartDialog,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"subs:date:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuStartDialog,
			h.adminOnly,
		),
	}
}

func (h *SubscriptionsTGHandler) Flows() []core_tg_dialog.Flow {
	return []core_tg_dialog.Flow{h.extendDaysFlow, h.setDateFlow}
}
