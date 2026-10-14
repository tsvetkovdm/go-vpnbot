package core_tg_dialog

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	CallbackPrefix = "dlg:"
	CancelValue    = "cancel"
)

func ButtonData(value string) string {
	return CallbackPrefix + value
}

func StepMarkup(step Step) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	if step.Markup != nil {
		rows = append(rows, step.Markup.InlineKeyboard...)
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", ButtonData(CancelValue)),
	))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
