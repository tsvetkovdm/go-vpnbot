package plans_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *PlansTGHandler) MenuViewPlan(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuViewPlan handler")

	data := update.CallbackQuery.Data
	planID, filter, page, err := parseCallbackData(data)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to parse callback data")
		return
	}

	plan, err := h.plansService.GetPlanByID(ctx, planID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get plan by ID")
		return
	}

	text, markup := assemblePlanCard(plan, filter, page)

	responseHandler.EditResponse(text, &markup)
}
