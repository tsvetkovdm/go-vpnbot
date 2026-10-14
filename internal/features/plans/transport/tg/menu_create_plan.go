package plans_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *PlansTGHandler) MenuCreatePlan(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)
	chatID := core_tg_request.ChatID(update)

	startStep := h.createPlanFlow.Steps[h.createPlanFlow.Start]
	responseHandler.EditResponse(
		startStep.Enter,
		nil,
	)

	session := core_tg_dialog.NewSession(
		h.createPlanFlow.ID,
		update.CallbackQuery.Message.MessageID,
		h.createPlanFlow.Start,
	)
	h.dialogStore.Set(chatID, session)
}
