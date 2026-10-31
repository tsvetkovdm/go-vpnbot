package plans_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *PlansTGHandler) MenuViewPlan(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuViewPlan handler")

	messageID := update.CallbackQuery.Message.MessageID

	planID, filter, page, err := parseCallbackData(update.CallbackQuery.Data)
	if err != nil {
		markup := plansMenuMarkup()
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &markup)
		return
	}

	text, markup, err := loadPlanCard(ctx, h.plansService, planID, filter, page)
	if err != nil {
		backMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(backToListButton(filter, page)))
		responseHandler.EditErrorResponse(messageID, err, "failed to load plan card", &backMarkup)
		return
	}

	responseHandler.EditResponse(text, &markup)
}
