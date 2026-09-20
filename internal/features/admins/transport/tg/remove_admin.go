package admins_transport_tg

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *AdminsTGHandler) RemoveAdmin(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke RemoveAdmin handler")

	telegramID, err := core_tg_request.TelegramIDArg(update)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to parse telegram_id argument")
		return
	}

	if err := h.adminsService.RemoveAdmin(ctx, telegramID); err != nil {
		responseHandler.ErrorResponse(err, "failed to remove admin")
		return
	}

	responseHandler.TextResponse(fmt.Sprintf("Администратор %d удалён.", telegramID))
}
