package admins_transport_tg

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
	"go.uber.org/zap"
)

func (h *AdminsTGHandler) MenuRoot(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)

	log.Debug("invoke MenuRoot admins handler")

	paymentsLabel := "💰 Платежи"
	count, err := h.adminsService.CountPendingPayments(ctx)
	if err != nil {
		log.Warn("failed to count pending payments", zap.Error(err))
	} else if count > 0 {
		paymentsLabel = fmt.Sprintf("💰 Платежи (%d)", count)
	}

	text := "Админ меню"
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(paymentsLabel, "pay:admin"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👥 Пользователи", "users:menu"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Тарифы", "plans:menu"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👤 Личный кабинет", "users:home"),
		),
	)

	if update.CallbackQuery != nil {
		core_tg_response.NewTGResponseHandler(log, update, h.bot).EditResponse(text, &markup)
		return
	}

	h.screens.Show(ctx, update, text, markup)
}
