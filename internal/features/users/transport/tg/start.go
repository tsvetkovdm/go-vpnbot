package users_transport_tg

import (
	"context"
	"time"

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
		userDomain = domain.NewUser(userTG.ID, nil, userTG.FirstName)
	} else {
		userDomain = domain.NewUser(userTG.ID, &userTG.UserName, userTG.FirstName)
	}

	user, err := h.usersService.UpsertUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to upsert user")
		return
	}

	overview, err := h.usersService.GetUserOverview(ctx, user.ID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user overview")
		return
	}

	text, markup := assembleHome(overview, time.Now())
	h.screens.Show(ctx, update, text, markup)
}
