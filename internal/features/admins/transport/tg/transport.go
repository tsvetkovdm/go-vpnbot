package admins_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_poller "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/poller"
)

type AdminsTGHandler struct {
	bot           *tgbotapi.BotAPI
	adminsService AdminsService
}

type AdminsService interface {
	GetAdminByTelegramID(
		ctx context.Context,
		telegramID int64,
	) (domain.Admin, error)
}

func NewAdminsTGHandler(
	bot *tgbotapi.BotAPI,
	adminsService AdminsService,
) *AdminsTGHandler {
	return &AdminsTGHandler{
		bot:           bot,
		adminsService: adminsService,
	}
}

func (h *AdminsTGHandler) Routes() []core_tg_poller.Route {
	return []core_tg_poller.Route{}
}
