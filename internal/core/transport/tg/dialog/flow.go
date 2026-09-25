package core_tg_dialog

type Flow struct {
	ID    string
	Start StepID
	Steps map[StepID]Step
}

func NewFlow(id string, start StepID, steps ...Step) Flow {
	m := make(map[StepID]Step, len(steps))
	for _, s := range steps {
		m[s.ID] = s
	}

	return Flow{
		ID:    id,
		Start: start,
		Steps: m,
	}
}
