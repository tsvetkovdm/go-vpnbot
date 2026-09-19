package users_transport_tg

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *UsersTGHandler) BlockUser(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke BlockUser handler")

	arg := update.Message.CommandArguments()
	telegramID, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("parse telegram_id %q: %w", arg, core_errors.ErrInvalidArgument),
			"failed to parse telegram_id argument",
		)

		return
	}

	if err := h.usersService.BlockUser(ctx, telegramID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to block user",
		)

		return
	}

	responseHandler.TextResponse(fmt.Sprintf("Пользователь %d заблокирован.", telegramID))
}
