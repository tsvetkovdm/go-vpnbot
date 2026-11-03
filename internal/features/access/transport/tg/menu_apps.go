package access_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *AccessTGHandler) MenuApps(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuApps handler")

	segments := core_tg_callback.Segments(update.CallbackQuery.Data)
	if len(segments) != 3 {
		h.MenuConnect(ctx, update)
		return
	}

	p, ok := findPlatform(segments[2])
	if !ok {
		h.MenuConnect(ctx, update)
		return
	}

	markup := appsMarkup(p)
	responseHandler.EditResponse(appsText(p), &markup)
}
