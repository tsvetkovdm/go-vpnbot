package core_tg_dialog

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StepID string

type Step struct {
	ID      StepID
	Enter   string
	Markup  *tgbotapi.InlineKeyboardMarkup
	OnInput func(ctx context.Context, session *Session, input string) Result
}
