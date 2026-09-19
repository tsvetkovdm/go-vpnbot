package users_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_middleware "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/middleware"
	core_tg_poller "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/poller"
)

type UsersTGHandler struct {
	bot          *tgbotapi.BotAPI
	usersService UsersService
	adminOnly    core_tg_middleware.Middleware
}

type UsersService interface {
	GetOrCreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	BlockUser(
		ctx context.Context,
		telegramID int64,
	) error

	UnblockUser(
		ctx context.Context,
		telegramID int64,
	) error

	DeleteUser(
		ctx context.Context,
		telegramID int64,
	) error
}

func NewUsersTGHandler(
	bot *tgbotapi.BotAPI,
	usersService UsersService,
	adminOnly core_tg_middleware.Middleware,
) *UsersTGHandler {
	return &UsersTGHandler{
		bot:          bot,
		usersService: usersService,
		adminOnly:    adminOnly,
	}
}

func (h *UsersTGHandler) Routes() []core_tg_poller.Route {
	return []core_tg_poller.Route{
		core_tg_poller.NewRoute(
			"start",
			core_tg_poller.TriggerCommand,
			h.Start,
		),
		core_tg_poller.NewRoute(
			"block",
			core_tg_poller.TriggerCommand,
			h.BlockUser,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"unblock",
			core_tg_poller.TriggerCommand,
			h.UnblockUser,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"delete",
			core_tg_poller.TriggerCommand,
			h.DeleteUser,
			h.adminOnly,
		),
	}
}
