package plans_transport_tg

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *PlansTGHandler) ListPlans(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke ListPlans handler")

	includeInactive := strings.EqualFold(
		strings.TrimSpace(update.Message.CommandArguments()),
		"all",
	)

	plans, err := h.plansService.ListPlans(ctx, includeInactive)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get list of plans")
		return
	}

	if len(plans) == 0 {
		responseHandler.TextResponse("Тарифов не найдено.")
		return
	}

	lines := make([]string, 0, len(plans))
	for _, plan := range plans {
		lines = append(lines, formatPlanLine(plan))
	}

	responseHandler.TextResponse(strings.Join(lines, "\n"))
}

func formatPlanLine(plan domain.Plan) string {
	status := "выключен"
	if plan.IsActive {
		status = "активен"
	}

	var duration string
	if plan.Duration.Months > 0 {
		duration = fmt.Sprintf("%d мес.", plan.Duration.Months)
	} else {
		duration = fmt.Sprintf("%d дн.", plan.Duration.Days)
	}

	return fmt.Sprintf(
		"#%d - %s - %s - %v р. - %s",
		plan.ID,
		plan.Title,
		duration,
		plan.Price,
		status,
	)
}
