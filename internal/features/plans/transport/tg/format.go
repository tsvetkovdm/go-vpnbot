package plans_transport_tg

import (
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func formatStatus(plan domain.Plan) string {
	status := "выключен"
	if plan.IsActive {
		status = "активен"
	}

	return status
}

func formatDuration(plan domain.Plan) string {
	var duration string
	if plan.Duration.Months > 0 {
		duration = fmt.Sprintf("%d мес.", plan.Duration.Months)
	} else {
		duration = fmt.Sprintf("%d дн.", plan.Duration.Days)
	}

	return duration
}

func formatPlanLine(plan domain.Plan) string {
	status := formatStatus(plan)
	duration := formatDuration(plan)

	return fmt.Sprintf(
		"#%d - %s - %s - %v р. - %s",
		plan.ID,
		plan.Title,
		duration,
		plan.Price,
		status,
	)
}

func formatPlanInfo(plan domain.Plan) string {
	status := formatStatus(plan)
	duration := formatDuration(plan)

	return fmt.Sprintf(
		"План #%d\nНазвание: %s\nДлительность: %s\nЦена: %v р.\nСтатус: %s\nСоздан: %s\nОбновлён: %s",
		plan.ID,
		plan.Title,
		duration,
		plan.Price,
		status,
		plan.CreatedAt.Format("2006-01-02 15:04"),
		plan.UpdatedAt.Format("2006-01-02 15:04"),
	)
}
