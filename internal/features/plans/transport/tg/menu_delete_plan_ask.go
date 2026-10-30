package plans_transport_tg

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *PlansTGHandler) MenuDeletePlanAsk(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuDeletePlanAsk handler")

	messageID := update.CallbackQuery.Message.MessageID

	planID, filter, page, err := parseCallbackData(update.CallbackQuery.Data)
	if err != nil {
		markup := plansMenuMarkup()
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &markup)
		return
	}

	// читаем тариф ради названия в вопросе и чтобы сразу поймать «уже удалён»
	plan, err := h.plansService.GetPlanByID(ctx, planID)
	if err != nil {
		markup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(backToListButton(filter, page)))
		responseHandler.EditErrorResponse(messageID, err, "failed to get plan", &markup)
		return
	}

	text := fmt.Sprintf(
		"Удалить тариф «%s»?\n\nУдалить можно только тариф, который ни разу не использовался. Действие нельзя отменить.",
		plan.Title,
	)
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Да, удалить",
				core_tg_callback.Build("plans", "confirm_delete", strconv.FormatInt(plan.ID, 10), filter, page),
			),
			viewPlanButton(plan.ID, filter, page, "Нет"),
		),
	)

	responseHandler.EditResponse(text, &markup)
}
