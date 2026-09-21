package core_tg_dialog

type ResultKind int

const (
	Retry ResultKind = iota
	Advance
	Complete
)

type Result struct {
	Kind    ResultKind
	Next    StepID // валиден только при Advance
	Message string // валиден при Retry (текст ошибки) и Complete (что ответить)
}
