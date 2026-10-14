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
	"go.uber.org/zap"
)

func (h *PlansTGHandler) MenuDeletePlanConfirm(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuDeletePlanConfirm handler")

	messageID := update.CallbackQuery.Message.MessageID

	planID, filter, page, err := parseCallbackData(update.CallbackQuery.Data)
	if err != nil {
		markup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад", core_tg_callback.Build("plans", "menu")),
			),
		)
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &markup)
		return
	}

	backMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(backToListButton(filter, page)))

	err = h.plansService.DeletePlan(ctx, planID)
	switch {
	case err == nil:
		responseHandler.EditResponse(fmt.Sprintf("Тариф #%d удалён.", planID), &backMarkup)

	case errors.Is(err, core_errors.ErrConflict):
		log.Info("delete plan rejected: plan in use", zap.Int64("plan_id", planID))
		markup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(viewPlanButton(planID, filter, page, "К тарифу")),
		)
		responseHandler.EditResponse("Тариф уже использовался — удалить нельзя, можно только выключить.", &markup)

	default:
		responseHandler.EditErrorResponse(messageID, err, "failed to delete plan", &backMarkup)
	}
}
