package core_tg_dialog

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type ResultKind int

const (
	Retry ResultKind = iota
	Advance
	Complete
	Abort
)

type Result struct {
	Kind ResultKind
	Next StepID // используется только при Advance

	// Message - зависит от Kind:
	//	 Retry: 	текст ошибки валидации, показывается как есть
	//	 Advance: 	необязательное переопределение Step.Enter следущего шага
	// 	 Complete: 	текст успеха, показывается как есть
	//	 Abort: 	не используется - текст формирует ErrorResponse из Err
	Message string

	// Err и LogMsg используются только при Abort
	Err    error
	LogMsg string

	// Markup — только для Complete: клавиатура вместо flow.ExitMarkup (nil — использовать ExitMarkup).
	Markup *tgbotapi.InlineKeyboardMarkup
}
