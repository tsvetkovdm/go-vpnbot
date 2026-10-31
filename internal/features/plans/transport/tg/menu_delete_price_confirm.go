package plans_transport_tg

import (
	"context"
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
	"go.uber.org/zap"
)

func (h *PlansTGHandler) MenuDeletePriceConfirm(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuDeletePriceConfirm handler")

	messageID := update.CallbackQuery.Message.MessageID

	priceID, filter, page, err := parseCallbackData(update.CallbackQuery.Data)
	if err != nil {
		markup := plansMenuMarkup()
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &markup)
		return
	}

	backMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(backToListButton(filter, page)))

	price, err := h.plansService.GetPlanPriceByID(ctx, priceID)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to get plan price", &backMarkup)
		return
	}

	var header string
	var text string
	var markup tgbotapi.InlineKeyboardMarkup

	err = h.plansService.DeletePlanPrice(ctx, priceID)
	switch {
	case err == nil:
		header = "✅ Срок удалён.\n\n"
		text, markup, err = loadPlanCard(ctx, h.plansService, price.PlanID, filter, page)

	case errors.Is(err, core_errors.ErrConflict):
		log.Info("delete plan price rejected: price in use", zap.Int64("price_id", priceID))
		header = "⚠️ Этот срок уже покупали — удалить нельзя, можно только выключить.\n\n"
		text, markup, err = loadPriceCard(ctx, h.plansService, priceID, filter, page)

	default:
		responseHandler.EditErrorResponse(messageID, err, "failed to delete plan price", &backMarkup)
		return
	}

	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to load card", &backMarkup)
		return
	}

	responseHandler.EditResponse(header+text, &markup)
}
