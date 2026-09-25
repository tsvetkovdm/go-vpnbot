package core_tg_dialog

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StepID string

type Step struct {
	ID      StepID
	Enter   string
	OnInput func(ctx context.Context, session *Session, update tgbotapi.Update) Result
}
