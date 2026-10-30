package plans_transport_tg

import (
	"context"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *PlansTGHandler) MenuTogglePlan(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuTogglePlan handler")

	messageID := update.CallbackQuery.Message.MessageID
	data := update.CallbackQuery.Data

	planID, filter, page, err := parseCallbackData(data)
	if err != nil {
		markup := plansMenuMarkup()
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &markup)
		return
	}

	backMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(backToListButton(filter, page)),
	)

	switch action := core_tg_callback.Segments(data)[1]; action {
	case "activate":
		_, err = h.plansService.ActivatePlan(ctx, planID)
	case "deactivate":
		_, err = h.plansService.DeactivatePlan(ctx, planID)
	default:
		err = fmt.Errorf("unknown plan action %q", action)
	}

	// ErrConflict = тариф уже в нужном состоянии (повторное нажатие) — просто показываем карточку
	if err != nil && !errors.Is(err, core_errors.ErrConflict) {
		responseHandler.EditErrorResponse(messageID, err, "failed to toggle plan", &backMarkup)
		return
	}

	text, markup, err := h.loadPlanCard(ctx, planID, filter, page)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to load plan card", &backMarkup)
		return
	}

	responseHandler.EditResponse(text, &markup)
}
