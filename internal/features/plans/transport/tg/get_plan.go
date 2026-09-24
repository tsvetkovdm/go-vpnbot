package plans_transport_tg

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *PlansTGHandler) GetPlan(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke GetPlan handler")

	planID, err := core_tg_request.Int64Arg(update)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get plan_id")
		return
	}

	plan, err := h.plansService.GetPlanByID(ctx, planID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get plan")
		return
	}

	responseHandler.TextResponse(formatPlanInfo(plan))
}

func formatPlanInfo(plan domain.Plan) string {
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
