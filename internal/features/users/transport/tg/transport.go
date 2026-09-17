package users_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_poller "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/poller"
)

type UsersTGHandler struct {
	bot          *tgbotapi.BotAPI
	usersService UsersService
}

type UsersService interface {
	GetOrCreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
}

func NewUsersTGHandler(
	bot *tgbotapi.BotAPI,
	usersService UsersService,
) *UsersTGHandler {
	return &UsersTGHandler{
		bot:          bot,
		usersService: usersService,
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
			"???",
			core_tg_poller.TriggerCallbackPrefix,
			h.Start,
		),
	}
}
