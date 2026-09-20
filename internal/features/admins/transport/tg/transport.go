package admins_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_middleware "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/middleware"
	core_tg_poller "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/poller"
)

type AdminsTGHandler struct {
	bot            *tgbotapi.BotAPI
	adminsService  AdminsService
	superAdminOnly core_tg_middleware.Middleware
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
}

func NewAdminsTGHandler(
	bot *tgbotapi.BotAPI,
	adminsService AdminsService,
	superAdminOnly core_tg_middleware.Middleware,
) *AdminsTGHandler {
	return &AdminsTGHandler{
		bot:            bot,
		adminsService:  adminsService,
		superAdminOnly: superAdminOnly,
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
	}
}
