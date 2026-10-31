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

func (h *PlansTGHandler) MenuPublicOffer(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuPublicOffer handler")

	messageID := update.CallbackQuery.Message.MessageID
	backMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ К тарифам", "plans:public"),
		),
	)

	fail := func(err error) {
		if errors.Is(err, core_errors.ErrNotFound) {
			log.Info("public plan unavailable", zap.Error(err))
			responseHandler.EditResponse("Этот тариф больше недоступен.", &backMarkup)
			return
		}
		responseHandler.EditErrorResponse(messageID, err, "failed to load public offer", &backMarkup)
	}

	planID, err := parseIDArg(update.CallbackQuery.Data)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &backMarkup)
		return
	}

	plan, err := h.plansService.GetPlanByID(ctx, planID)
	if err != nil {
		fail(err)
		return
	}

	text, markup, err := loadPublicOffer(ctx, h.plansService, plan, "plans:public")
	if err != nil {
		fail(err)
		return
	}

	responseHandler.EditResponse(text, &markup)
}
