package users_transport_tg

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

func (h *UsersTGHandler) DeleteUser(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke DeleteUser handler")

	telegramID, err := core_tg_request.TelegramIDArg(update)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to parse telegram_id argument")
		return
	}

	if err := h.usersService.DeleteUser(ctx, telegramID); err != nil {
		if errors.Is(err, core_errors.ErrConflict) {
			responseHandler.TextResponse(fmt.Sprintf(
				"У пользователя %d есть подписка или история платежей - "+
					"удалить нельзя, только заблокировать (/block).",
				telegramID,
			))
			return
		}

		responseHandler.ErrorResponse(err, "failed to delete user")
		return
	}

	responseHandler.TextResponse(fmt.Sprintf("Пользователь %d удалён.", telegramID))
}
