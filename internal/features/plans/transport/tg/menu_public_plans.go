package plans_transport_tg

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
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
	backMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "users:home"),
		),
	)

	plans, err := h.plansService.ListPlans(ctx, false)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to list public plans", &backMarkup)
		return
	}

	if len(plans) == 0 {
		responseHandler.EditResponse("💳 Тарифы\n\nСейчас нет доступных тарифов.", &backMarkup)
		return
	}

	var sb strings.Builder
	sb.WriteString("💳 Тарифы\n\n")
	for _, plan := range plans {
		sb.WriteString(fmt.Sprintf("• %s\n", plan.Title))
	}
	sb.WriteString("\nДля оплаты напишите администратору.")

	responseHandler.EditResponse(sb.String(), &backMarkup)
}
