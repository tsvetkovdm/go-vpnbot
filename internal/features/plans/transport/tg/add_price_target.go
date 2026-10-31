package plans_transport_tg

import (
	"fmt"
	"strconv"

	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
)

const (
	dataKeyPlanID = "plan_id"
	dataKeyFilter = "filter"
	dataKeyPage   = "page"
)

// planTarget — тариф, к которому добавляем срок, и куда вернуться (filter/page списка)
type planTarget struct {
	planID int64
	filter string
	page   string
}

func (t planTarget) saveTo(session *core_tg_dialog.Session) {
	session.Data[dataKeyPlanID] = strconv.FormatInt(t.planID, 10)
	session.Data[dataKeyFilter] = t.filter
	session.Data[dataKeyPage] = t.page
}

func loadPlanTarget(session *core_tg_dialog.Session) (planTarget, error) {
	planID, err := strconv.ParseInt(session.Data[dataKeyPlanID], 10, 64)
	if err != nil {
		return planTarget{}, fmt.Errorf("session plan_id: %w", err)
	}

	return planTarget{
		planID: planID,
		filter: session.Data[dataKeyFilter],
		page:   session.Data[dataKeyPage],
	}, nil
}
