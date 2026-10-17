package users_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
	core_tg_middleware "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/middleware"
	core_tg_poller "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/poller"
	core_tg_screen "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/screen"
)

type UsersTGHandler struct {
	bot            *tgbotapi.BotAPI
	usersService   UsersService
	adminOnly      core_tg_middleware.Middleware
	dialogStore    *core_tg_dialog.Store
	searchUserFlow core_tg_dialog.Flow
	screens        *core_tg_screen.Screens
}

type UsersService interface {
	UpsertUser(
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

	GetUserOverview(
		ctx context.Context,
		userID int64,
	) (domain.UserOverview, error)

	GetUserOverviewByTelegramID(
		ctx context.Context,
		telegramID int64,
	) (domain.UserOverview, error)

	GetUserOverviewByUsername(
		ctx context.Context,
		username string,
	) (domain.UserOverview, error)

	ListUserOverviews(
		ctx context.Context,
		filter domain.UserFilter,
		limit, offset int,
	) ([]domain.UserOverview, int, error)
}

func NewUsersTGHandler(
	bot *tgbotapi.BotAPI,
	usersService UsersService,
	adminOnly core_tg_middleware.Middleware,
	dialogStore *core_tg_dialog.Store,
	screens *core_tg_screen.Screens,
) *UsersTGHandler {
	return &UsersTGHandler{
		bot:            bot,
		usersService:   usersService,
		adminOnly:      adminOnly,
		dialogStore:    dialogStore,
		searchUserFlow: newSearchUserFlow(usersService),
		screens:        screens,
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
		core_tg_poller.NewRoute(
			"users:home",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuHome,
		),
		core_tg_poller.NewRoute(
			"users:menu",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuRoot,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"users:list:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuListUsers,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"users:view:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuViewUser,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"users:block:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuBlockUser,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"users:unblock:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuBlockUser,
			h.adminOnly,
		),
		core_tg_poller.NewRoute(
			"users:search",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuSearchUser,
			h.adminOnly,
		),
	}
}

func (h *UsersTGHandler) Flows() []core_tg_dialog.Flow {
	return []core_tg_dialog.Flow{h.searchUserFlow}
}
