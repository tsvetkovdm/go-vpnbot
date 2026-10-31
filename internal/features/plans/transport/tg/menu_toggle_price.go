package plans_transport_tg

import (
	"context"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
	"go.uber.org/zap"
)

func (h *PlansTGHandler) MenuTogglePrice(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuTogglePrice handler")

	messageID := update.CallbackQuery.Message.MessageID
	data := update.CallbackQuery.Data

	priceID, filter, page, err := parseCallbackData(data)
	if err != nil {
		markup := plansMenuMarkup()
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &markup)
		return
	}

	backMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(backToListButton(filter, page)),
	)

	switch action := core_tg_callback.Segments(data)[1]; action {
	case "pactivate":
		_, err = h.plansService.ActivatePlanPrice(ctx, priceID)
	case "pdeactivate":
		_, err = h.plansService.DeactivatePlanPrice(ctx, priceID)
	default:
		err = fmt.Errorf("unknown price action %q", action)
	}

	var header string
	switch {
	case err == nil:
	case errors.Is(err, domain.ErrPlanPriceDuplicateDuration):
		// проверять раньше ErrConflict: дубль — тоже ErrConflict
		log.Info("activate plan price rejected: duplicate duration", zap.Error(err))
		header = "⚠️ Не включён: в тарифе уже есть активный срок такой же длительности " +
			"(1 месяц и 30 дней считаются одинаковыми). Сначала выключите его.\n\n"
	case errors.Is(err, core_errors.ErrConflict):
		// срок уже в нужном состоянии (повторное нажатие) — просто показываем карточку
	default:
		responseHandler.EditErrorResponse(messageID, err, "failed to toggle plan price", &backMarkup)
		return
	}

	text, markup, err := loadPriceCard(ctx, h.plansService, priceID, filter, page)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to load price card", &backMarkup)
		return
	}

	responseHandler.EditResponse(header+text, &markup)
}
