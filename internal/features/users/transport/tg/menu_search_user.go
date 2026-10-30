package users_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *UsersTGHandler) MenuSearchUser(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)
	chatID := core_tg_request.ChatID(update)

	log.Debug("invoke MenuSearchUser handler")

	startStep := h.searchUserFlow.Steps[h.searchUserFlow.Start]
	markup := core_tg_dialog.StepMarkup(startStep)
	responseHandler.EditResponse(startStep.Enter, &markup)

	session := core_tg_dialog.NewSession(
		h.searchUserFlow.ID,
		update.CallbackQuery.Message.MessageID,
		h.searchUserFlow.Start,
	)
	h.dialogStore.Set(chatID, session)
}
