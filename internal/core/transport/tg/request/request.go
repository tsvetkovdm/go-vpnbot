package core_tg_request

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
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

func TelegramIDArg(update tgbotapi.Update) (int64, error) {
	arg := update.Message.CommandArguments()

	telegramID, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse telegram_id %q: %w", arg, core_errors.ErrInvalidArgument)
	}

	return telegramID, nil
}
