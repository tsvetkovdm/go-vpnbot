package plans_transport_tg

import (
	"context"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *PlansTGHandler) MenuPublicPlans(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuPublicPlans handler")

	messageID := update.CallbackQuery.Message.MessageID
	homeMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "users:home"),
		),
	)

	plans, err := h.plansService.ListPublicPlans(ctx)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to list public plans", &homeMarkup)
		return
	}

	switch len(plans) {
	case 0:
		responseHandler.EditResponse("💳 Тарифы\n\nСейчас нет доступных тарифов.", &homeMarkup)

	case 1:
		text, markup, err := loadPublicOffer(ctx, h.plansService, plans[0], "users:home")
		if err != nil {
			responseHandler.EditErrorResponse(messageID, err, "failed to load public offer", &homeMarkup)
			return
		}
		responseHandler.EditResponse(text, &markup)

	default:
		rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(plans)+1)
		for _, plan := range plans {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					plan.Title,
					core_tg_callback.Build("plans", "offer", strconv.FormatInt(plan.ID, 10)),
				),
			))
		}
		rows = append(rows, homeMarkup.InlineKeyboard...)

		markup := tgbotapi.NewInlineKeyboardMarkup(rows...)
		responseHandler.EditResponse("💳 Тарифы\n\nВыберите тариф:", &markup)
	}
}
