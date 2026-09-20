package core_tg_middleware

import (
	"context"
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_handler "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/handler"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

type AdminsService interface {
	GetAdminByTelegramID(
		ctx context.Context,
		telegramID int64,
	) (domain.Admin, error)
}

func AdminOnly(bot *tgbotapi.BotAPI, adminsService AdminsService) Middleware {
	return func(next core_tg_handler.Handler) core_tg_handler.Handler {
		return func(ctx context.Context, update tgbotapi.Update) {
			log := core_logger.FromContext(ctx)
			responseHandler := core_tg_response.NewTGResponseHandler(log, update, bot)

			user := core_tg_request.FromUser(update)
			_, ok, err := checkAdmin(ctx, adminsService, user.ID)
			if err != nil {
				responseHandler.ErrorResponse(err, "failed to check admin access")
				return
			}

			if !ok {
				return
			}

			next(ctx, update)
		}
	}
}

func SuperAdminOnly(bot *tgbotapi.BotAPI, adminsService AdminsService) Middleware {
	return func(next core_tg_handler.Handler) core_tg_handler.Handler {
		return func(ctx context.Context, update tgbotapi.Update) {
			log := core_logger.FromContext(ctx)
			responseHandler := core_tg_response.NewTGResponseHandler(log, update, bot)

			user := core_tg_request.FromUser(update)

			admin, ok, err := checkAdmin(ctx, adminsService, user.ID)
			if err != nil {
				responseHandler.ErrorResponse(err, "failed to check admin access")
				return
			}

			if !ok || admin.Role != domain.RoleSuper {
				return
			}

			next(ctx, update)
		}
	}
}

func checkAdmin(
	ctx context.Context,
	adminsService AdminsService,
	telegramID int64,
) (domain.Admin, bool, error) {
	admin, err := adminsService.GetAdminByTelegramID(ctx, telegramID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return domain.Admin{}, false, nil
		}
		return domain.Admin{}, false, err
	}

	return admin, true, nil
}
