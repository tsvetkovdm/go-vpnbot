package plans_transport_tg

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
	"go.uber.org/zap"
)

func (h *PlansTGHandler) MenuPublicBuy(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuPublicBuy handler")

	messageID := update.CallbackQuery.Message.MessageID
	backRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ К тарифам", "plans:public"),
		tgbotapi.NewInlineKeyboardButtonData("👤 Личный кабинет", "users:home"),
	)
	backMarkup := tgbotapi.NewInlineKeyboardMarkup(backRow)

	fail := func(err error) {
		if errors.Is(err, core_errors.ErrNotFound) {
			log.Info("public plan price unavailable", zap.Error(err))
			responseHandler.EditResponse("Этот вариант больше недоступен.", &backMarkup)
			return
		}
		responseHandler.EditErrorResponse(messageID, err, "failed to load plan price", &backMarkup)
	}

	priceID, err := parseIDArg(update.CallbackQuery.Data)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &backMarkup)
		return
	}

	price, err := h.plansService.GetPlanPriceByID(ctx, priceID)
	if err != nil {
		fail(err)
		return
	}
	if !price.IsActive {
		fail(fmt.Errorf("plan price id=%d is inactive: %w", price.ID, errPlanUnavailable))
		return
	}

	plan, err := h.plansService.GetPlanByID(ctx, price.PlanID)
	if err != nil {
		fail(err)
		return
	}
	if !plan.IsActive {
		fail(fmt.Errorf("plan id=%d is inactive: %w", plan.ID, errPlanUnavailable))
		return
	}

	if price.Price.IsZero() {
		h.showFreePrice(ctx, responseHandler, update, plan, price, backRow)
		return
	}

	text := fmt.Sprintf(
		"🧾 Ваш выбор\n\nТариф: %s\nСрок: %s\nЦена: %s\n\n"+
			"Нажмите «💳\u00a0Оплатить», чтобы получить реквизиты для перевода.",
		plan.Title,
		formatPriceDuration(price.Duration),
		formatPriceAmount(price),
	)

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"💳 Оплатить",
				core_tg_callback.Build("pay", "create", strconv.FormatInt(price.ID, 10)),
			),
		),
		backRow,
	)

	responseHandler.EditResponse(text, &markup)
}
