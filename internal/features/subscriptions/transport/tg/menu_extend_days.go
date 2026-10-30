package subscriptions_transport_tg

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
	"go.uber.org/zap"
)

type extendDaysCallback struct {
	userID  int64
	days    int
	version int64
	filter  string
	page    string
}

func (h *SubscriptionsTGHandler) MenuExtendDays(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuExtendDays handler")

	messageID := update.CallbackQuery.Message.MessageID

	cb, err := parseExtendDaysCallback(update.CallbackQuery.Data)
	if err != nil {
		adminMenu := adminMenuMarkup()
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &adminMenu)
		return
	}

	backMarkup := backToCardMarkup(cb.userID, cb.filter, cb.page)

	sub, err := h.subscriptionsService.ExtendDays(ctx, cb.userID, cb.days, cb.version)

	if errors.Is(err, core_errors.ErrConflict) {
		log.Info("subscription changed concurrently", zap.Error(err))
		responseHandler.EditResponse(
			conflictText,
			&backMarkup,
		)
		return
	}

	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to extend subscription", &backMarkup)
		return
	}

	text := fmt.Sprintf(
		"✅ Доступ продлён на %d дн.\nОплачено до %s",
		cb.days,
		sub.PaidUntil(time.Local).Format("02.01.2006"),
	)
	responseHandler.EditResponse(text, &backMarkup)
}

func parseExtendDaysCallback(data string) (extendDaysCallback, error) {
	segments := core_tg_callback.Segments(data)
	if len(segments) != 7 {
		return extendDaysCallback{}, fmt.Errorf("malformed callback data %q: %w", data, core_errors.ErrInvalidArgument)
	}

	userID, err := strconv.ParseInt(segments[2], 10, 64)
	if err != nil {
		return extendDaysCallback{}, fmt.Errorf("invalid userID %q: %w", segments[2], core_errors.ErrInvalidArgument)
	}

	days, err := strconv.Atoi(segments[3])
	if err != nil {
		return extendDaysCallback{}, fmt.Errorf("invalid days %q: %w", segments[3], core_errors.ErrInvalidArgument)
	}

	version, err := strconv.ParseInt(segments[4], 10, 64)
	if err != nil {
		return extendDaysCallback{}, fmt.Errorf("invalid version %q: %w", segments[4], core_errors.ErrInvalidArgument)
	}

	return extendDaysCallback{
		userID:  userID,
		days:    days,
		version: version,
		filter:  segments[5],
		page:    segments[6],
	}, nil
}
