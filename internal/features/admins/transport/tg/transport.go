package admins_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_middleware "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/middleware"
	core_tg_poller "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/poller"
	core_tg_screen "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/screen"
)

type AdminsTGHandler struct {
	bot            *tgbotapi.BotAPI
	adminsService  AdminsService
	adminOnly      core_tg_middleware.Middleware
	superAdminOnly core_tg_middleware.Middleware
	screens        *core_tg_screen.Screens
}

type AdminsService interface {
	GetAdminByTelegramID(
		ctx context.Context,
		telegramID int64,
	) (domain.Admin, error)

	AddAdmin(
		ctx context.Context,
		admin domain.Admin,
	) (domain.Admin, error)

	RemoveAdmin(
		ctx context.Context,
		telegramID int64,
	) error

	CountPendingPayments(
		ctx context.Context,
	) (int, error)
}

func NewAdminsTGHandler(
	bot *tgbotapi.BotAPI,
	adminsService AdminsService,
	adminOnly core_tg_middleware.Middleware,
	superAdminOnly core_tg_middleware.Middleware,
	screens *core_tg_screen.Screens,
) *AdminsTGHandler {
	return &AdminsTGHandler{
		bot:            bot,
		adminsService:  adminsService,
		adminOnly:      adminOnly,
		superAdminOnly: superAdminOnly,
		screens:        screens,
	}
}

func (h *AdminsTGHandler) Routes() []core_tg_poller.Route {
	return []core_tg_poller.Route{
		core_tg_poller.NewRoute(
			"add_admin",
			core_tg_poller.TriggerCommand,
			h.AddAdmin,
			h.superAdminOnly,
		),
		core_tg_poller.NewRoute(
			"remove_admin",
			core_tg_poller.TriggerCommand,
			h.RemoveAdmin,
			h.superAdminOnly,
		),
		core_tg_poller.NewRoute(
			"admin",
			core_tg_poller.TriggerCommand,
			h.MenuRoot,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"admins:menu",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuRoot,
			h.adminOnly,
		),
	}
}
