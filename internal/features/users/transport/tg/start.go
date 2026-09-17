package users_transport_tg

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *UsersTGHandler) Start(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke Start handler")

	userTG := core_tg_request.FromUser(update)
	var userDomain domain.User
	if userTG.UserName == "" {
		userDomain = domain.NewUser(userTG.ID, nil)
	} else {
		userDomain = domain.NewUser(userTG.ID, &userTG.UserName)
	}

	userDomain, err := h.usersService.GetOrCreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create user",
		)
		return
	}

	responseHandler.TextResponse(mainMenuText(userDomain))
}

func mainMenuText(user domain.User) string {
	return fmt.Sprintf(
		`
		Добро пожаловать, %s!
		TelegramID: %d
		Email: %s
		XrayUUID: %s
		Status: %s
		CreatedAt: %s
		UpdatedAt: %s
		`,
		*user.Username,
		user.TelegramID,
		user.Email,
		user.XrayUUID,
		user.Status,
		user.CreatedAt,
		user.UpdatedAt,
	)
}
