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

func (h *PlansTGHandler) MenuDeletePriceAsk(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuDeletePriceAsk handler")

	messageID := update.CallbackQuery.Message.MessageID

	priceID, filter, page, err := parseCallbackData(update.CallbackQuery.Data)
	if err != nil {
		markup := plansMenuMarkup()
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &markup)
		return
	}

	price, err := h.plansService.GetPlanPriceByID(ctx, priceID)
	if err != nil {
		markup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(backToListButton(filter, page)))
		responseHandler.EditErrorResponse(messageID, err, "failed to get plan price", &markup)
		return
	}

	text := fmt.Sprintf(
		"Удалить срок «%s — %s»?\n\nУдалить можно только срок, который ни разу не покупали. Действие нельзя отменить.",
		formatPriceDuration(price.Duration),
		formatPriceAmount(price),
	)

	id := strconv.FormatInt(price.ID, 10)
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Да, удалить",
				core_tg_callback.Build("plans", "pconfirm_delete", id, filter, page),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"Нет",
				core_tg_callback.Build("plans", "pview", id, filter, page),
			),
		),
	)

	responseHandler.EditResponse(text, &markup)
}
