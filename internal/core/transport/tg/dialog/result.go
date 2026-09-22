package core_tg_dialog

type ResultKind int

const (
	Retry ResultKind = iota
	Advance
	Complete
)

type Result struct {
	Kind    ResultKind
	Next    StepID
	Message string
}
