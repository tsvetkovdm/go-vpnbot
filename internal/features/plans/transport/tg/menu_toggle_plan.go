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
		markup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					"Назад",
					core_tg_callback.Build("plans", "menu"),
				),
			),
		)
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &markup)
		return
	}

	backMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(backToListButton(filter, page)),
	)

	var plan domain.Plan
	switch action := core_tg_callback.Segments(data)[1]; action {
	case "activate":
		plan, err = h.plansService.ActivatePlan(ctx, planID)
	case "deactivate":
		plan, err = h.plansService.DeactivatePlan(ctx, planID)
	default:
		err = fmt.Errorf("unknown plan action %q", action)
	}

	if errors.Is(err, core_errors.ErrConflict) {
		plan, err = h.plansService.GetPlanByID(ctx, planID)
	}
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to toggle plan", &backMarkup)
		return
	}

	text, markup := assemblePlanCard(plan, filter, page)
	responseHandler.EditResponse(text, &markup)
}
