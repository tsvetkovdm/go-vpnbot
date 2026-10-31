package plans_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	text, _, err := loadPlanCard(ctx, h.plansService, planID, "all", "0")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to load plan card")
		return
	}

	responseHandler.TextResponse(text)
}
