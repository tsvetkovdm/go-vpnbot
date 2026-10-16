package subscriptions_transport_tg

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_middleware "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/middleware"
	core_tg_poller "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/poller"
)

type SubscriptionsTGHandler struct {
	bot                  *tgbotapi.BotAPI
	subscriptionsService SubscriptionsService
	adminOnly            core_tg_middleware.Middleware
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
) *SubscriptionsTGHandler {
	return &SubscriptionsTGHandler{
		bot:                  bot,
		subscriptionsService: subscriptionsService,
		adminOnly:            adminOnly,
	}
}

func (h *SubscriptionsTGHandler) Routes() []core_tg_poller.Route {
	return []core_tg_poller.Route{
		core_tg_poller.NewRoute(
			"subs:extend:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuExtendDays,
			h.adminOnly,
		),
	}
}
