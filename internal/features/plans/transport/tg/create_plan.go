package plans_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
)

func (h *PlansTGHandler) CreatePlan(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)

	log.Debug("invoke CreatePlan handler")

	chatID := core_tg_request.ChatID(update)
	startStep := h.createPlanFlow.Steps[h.createPlanFlow.Start]

	messageID := h.screens.Show(ctx, update, startStep.Enter, core_tg_dialog.StepMarkup(startStep))
	if messageID == 0 {
		return
	}

	session := core_tg_dialog.NewSession(h.createPlanFlow.ID, messageID, h.createPlanFlow.Start)
	h.dialogStore.Set(chatID, session)
}
