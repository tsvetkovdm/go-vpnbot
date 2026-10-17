package users_transport_tg

import (
	"context"
	"errors"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *UsersTGHandler) MenuHome(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuHome handler")

	messageID := update.CallbackQuery.Message.MessageID
	telegramID := core_tg_request.FromUser(update).ID

	overview, err := h.usersService.GetUserOverviewByTelegramID(ctx, telegramID)
	if errors.Is(err, core_errors.ErrNotFound) {
		responseHandler.EditResponse("Нажмите /start, чтобы начать.", nil)
		return
	}
	if err != nil {
		refreshMarkup := tgbotapi.NewInlineKeyboardMarkup(refreshRow())
		responseHandler.EditErrorResponse(messageID, err, "failed to get user overview", &refreshMarkup)
		return
	}

	text, markup := assembleHome(overview, time.Now())
	responseHandler.EditResponse(text, &markup)
}
