package core_tg_dialog

import "time"

type Session struct {
	FlowID    string
	Step      StepID
	Data      map[string]string
	UpdatedAt time.Time
}

func NewSession(flowID string, step StepID) Session {
	return Session{
		FlowID:    flowID,
		Step:      step,
		Data:      make(map[string]string),
		UpdatedAt: time.Now(),
	}
}
