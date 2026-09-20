package admins_transport_tg

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *AdminsTGHandler) AddAdmin(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke AddAdmin handler")

	telegramID, err := core_tg_request.TelegramIDArg(update)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to parse telegram_id argument")
		return
	}

	adminTG := core_tg_request.FromUser(update)
	adminFrom, err := h.adminsService.GetAdminByTelegramID(ctx, adminTG.ID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get admin")
		return
	}

	adminDomain := domain.NewAdmin(
		telegramID,
		domain.RoleAdmin,
		&adminFrom.ID,
	)

	adminDomain, err = h.adminsService.AddAdmin(ctx, adminDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to add admin")
		return
	}

	responseHandler.TextResponse(fmt.Sprintf("Новый администратор %d добавлен. Роль: %v", telegramID, adminDomain.Role))
}
