package core_tg_dialog

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Flow struct {
	ID         string
	Start      StepID
	Steps      map[StepID]Step
	CancelText string
	ExitMarkup *tgbotapi.InlineKeyboardMarkup
}

func NewFlow(
	id string,
	start StepID,
	cancelText string,
	exitMarkup *tgbotapi.InlineKeyboardMarkup,
	steps ...Step,
) Flow {
	m := make(map[StepID]Step, len(steps))
	for _, s := range steps {
		m[s.ID] = s
	}

	return Flow{
		ID:         id,
		Start:      start,
		Steps:      m,
		CancelText: cancelText,
		ExitMarkup: exitMarkup,
	}
}
