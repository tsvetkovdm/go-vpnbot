package subscriptions_transport_tg

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
)

func backToCardMarkup(userID int64, filter, page string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"⬅️ К карточке",
				core_tg_callback.Build("users", "view", strconv.FormatInt(userID, 10), filter, page),
			),
		),
	)
}
