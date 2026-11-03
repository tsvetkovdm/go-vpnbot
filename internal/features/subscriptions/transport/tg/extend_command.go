package subscriptions_transport_tg

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *SubscriptionsTGHandler) ExtendCommand(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke ExtendCommand handler")

	args := strings.Fields(update.Message.CommandArguments())
	if len(args) != 2 {
		responseHandler.TextResponse("Использование: /extend <telegram_id> <дни>\nПример: /extend 1234567890 30")
		return
	}

	telegramID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		responseHandler.TextResponse("Telegram ID должен быть числом.")
		return
	}

	days, err := strconv.Atoi(args[1])
	if err != nil {
		responseHandler.TextResponse("Число дней должно быть целым числом.")
		return
	}

	userID, err := h.subscriptionsService.ResolveUserID(ctx, telegramID)
	if errors.Is(err, core_errors.ErrNotFound) {
		responseHandler.TextResponse(fmt.Sprintf(
			"Пользователь %d не найден. Он должен сначала нажать /start в боте.", telegramID,
		))
		return
	}
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to resolve user")
		return
	}

	sub, err := h.subscriptionsService.ExtendDays(ctx, userID, days, domain.SubscriptionSkipVersionCheck)
	if errors.Is(err, core_errors.ErrInvalidArgument) {
		responseHandler.TextResponse(fmt.Sprintf("Число дней должно быть от 1 до %d.", domain.SubscriptionMaxExtendDays))
		return
	}
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to extend subscription")
		return
	}

	responseHandler.TextResponse(fmt.Sprintf(
		"✅ %d: доступ продлён на %d дн.\nОплачено до %s",
		telegramID,
		days,
		sub.PaidUntil(time.Local).Format("02.01.2006"),
	))
}
