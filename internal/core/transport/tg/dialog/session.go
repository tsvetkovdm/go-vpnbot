package core_tg_dialog

import "time"

type Session struct {
	FlowID    string
	MessageID int
	Step      StepID
	Data      map[string]string
	UpdatedAt time.Time
}

func NewSession(flowID string, messageID int, step StepID) Session {
	return Session{
		FlowID:    flowID,
		MessageID: messageID,
		Step:      step,
		Data:      make(map[string]string),
		UpdatedAt: time.Now(),
	}
}
