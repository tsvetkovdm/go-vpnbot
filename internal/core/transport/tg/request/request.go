package core_tg_request

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ChatID(update tgbotapi.Update) int64 {
	var chatID int64

	switch {
	case update.Message != nil:
		chatID = update.Message.Chat.ID
		return chatID
	case update.CallbackQuery != nil:
		chatID = update.CallbackQuery.Message.Chat.ID
		return chatID
	default:
		panic(fmt.Sprintf(
			"core_tg_request: ChatID: update has neither Message nor CallbackQuery, update_id=%d",
			update.UpdateID,
		))
	}
}

func FromUser(update tgbotapi.Update) *tgbotapi.User {
	switch {
	case update.Message != nil:
		return update.Message.From
	case update.CallbackQuery != nil:
		return update.CallbackQuery.From
	default:
		panic(fmt.Sprintf(
			"core_tg_request: FromUser: update has neither Message nor CallbackQuery, update_id=%d",
			update.UpdateID,
		))

	}
}
