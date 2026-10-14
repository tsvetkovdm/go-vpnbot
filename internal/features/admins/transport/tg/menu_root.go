package admins_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *AdminsTGHandler) MenuRoot(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)

	log.Debug("invoke MenuRoot admins handler")

	text := "Админ меню"
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Тарифы", "plans:menu"),
		),
	)

	if update.CallbackQuery != nil {
		core_tg_response.NewTGResponseHandler(log, update, h.bot).EditResponse(text, &markup)
		return
	}

	h.screens.Show(ctx, update, text, markup)
}
