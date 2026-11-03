package access_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_tg_poller "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/poller"
)

type AccessTGHandler struct {
	bot           *tgbotapi.BotAPI
	accessService AccessService
	config        Config
}

type AccessService interface {
	GetFeedToken(
		ctx context.Context,
		telegramID int64,
	) (string, error)
}

func NewAccessTGHandler(
	bot *tgbotapi.BotAPI,
	accessService AccessService,
	config Config,
) *AccessTGHandler {
	return &AccessTGHandler{
		bot:           bot,
		accessService: accessService,
		config:        config,
	}
}

func (h *AccessTGHandler) Routes() []core_tg_poller.Route {
	return []core_tg_poller.Route{
		core_tg_poller.NewRoute(
			"access:connect",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuConnect,
		),
		core_tg_poller.NewRoute(
			"access:apps:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuApps,
		),
	}
}
