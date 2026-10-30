package users_transport_tg

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *UsersTGHandler) MenuBlockUser(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuBlockUser handler")

	messageID := update.CallbackQuery.Message.MessageID
	data := update.CallbackQuery.Data
	backMarkup := backToUsersMenuMarkup()

	userID, filter, page, err := parseUserCallback(data)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &backMarkup)
		return
	}

	backToList := backToListMarkup(filter, page)

	overview, err := h.usersService.GetUserOverview(ctx, userID)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to get user overview", &backToList)
		return
	}

	switch core_tg_callback.Segments(data)[1] {
	case "block":
		if err := h.usersService.BlockUser(ctx, overview.User.TelegramID); err != nil {
			responseHandler.EditErrorResponse(messageID, err, "failed to block user", &backToList)
			return
		}
	case "unblock":
		if err := h.usersService.UnblockUser(ctx, overview.User.TelegramID); err != nil {
			responseHandler.EditErrorResponse(messageID, err, "failed to unblock user", &backToList)
			return
		}
	}

	overview, err = h.usersService.GetUserOverview(ctx, userID)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to get user overview", &backToList)
		return
	}

	text, markup := assembleUserCard(overview, filter, page, time.Now())
	responseHandler.EditResponse(text, &markup)
}
