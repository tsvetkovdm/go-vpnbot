package plans_transport_tg

import (
	"context"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *PlansTGHandler) ActivatePlan(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke ActivatePlan handler")

	planID, err := core_tg_request.Int64Arg(update)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to parse plan_id")
		return
	}

	plan, err := h.plansService.ActivatePlan(ctx, planID)
	if err != nil {
		if errors.Is(err, core_errors.ErrConflict) {
			responseHandler.TextResponse(fmt.Sprintf("Тариф #%d уже включён.", planID))
			return
		}
		responseHandler.ErrorResponse(err, "failed to activate plan")
		return
	}

	responseHandler.TextResponse("Тариф включён: " + formatPlanLine(plan))
}
