package plans_transport_tg

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func (h *PlansTGHandler) showFreePrice(
	ctx context.Context,
	responseHandler *core_tg_response.TGResponseHandler,
	update tgbotapi.Update,
	plan domain.Plan,
	price domain.PlanPrice,
	backRow []tgbotapi.InlineKeyboardButton,
) {
	backMarkup := tgbotapi.NewInlineKeyboardMarkup(backRow)
	messageID := update.CallbackQuery.Message.MessageID

	canUse, err := h.plansService.CanUseFreePrices(ctx, core_tg_request.FromUser(update).ID)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to check free prices", &backMarkup)
		return
	}
	if !canUse {
		responseHandler.EditResponse("🎁 Бесплатный доступ — только для новых пользователей.", &backMarkup)
		return
	}

	text := fmt.Sprintf(
		"🎁 Бесплатный доступ\n\nТариф: %s\nСрок: %s\n\n"+
			"Попробуйте сервис бесплатно — оплачивать ничего не нужно.",
		plan.Title,
		formatPriceDuration(price.Duration),
	)

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🎁 Получить",
				core_tg_callback.Build("pay", "free", strconv.FormatInt(price.ID, 10)),
			),
		),
		backRow,
	)

	responseHandler.EditResponse(text, &markup)
}
