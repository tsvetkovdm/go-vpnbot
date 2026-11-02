package admins_transport_tg

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
	"go.uber.org/zap"
)

func (h *AdminsTGHandler) MenuRoot(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)

	log.Debug("invoke MenuRoot admins handler")

	var rows [][]tgbotapi.InlineKeyboardButton
	if h.isSuperAdmin(ctx, core_tg_request.FromUser(update).ID) {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(h.paymentsLabel(ctx), "pay:admin"),
		))
	}

	text := "Админ меню"
	rows = append(rows,
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
	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	if update.CallbackQuery != nil {
		core_tg_response.NewTGResponseHandler(log, update, h.bot).EditResponse(text, &markup)
		return
	}

	h.screens.Show(ctx, update, text, markup)
}

func (h *AdminsTGHandler) isSuperAdmin(ctx context.Context, telegramID int64) bool {
	admin, err := h.adminsService.GetAdminByTelegramID(ctx, telegramID)
	if err != nil {
		core_logger.FromContext(ctx).Warn("failed to get admin role", zap.Error(err))
		return false
	}

	return admin.Role == domain.RoleSuper
}

func (h *AdminsTGHandler) paymentsLabel(ctx context.Context) string {
	count, err := h.adminsService.CountPendingPayments(ctx)
	if err != nil {
		core_logger.FromContext(ctx).Warn("failed to count pending payments", zap.Error(err))
		return "💰 Платежи"
	}

	if count == 0 {
		return "💰 Платежи"
	}

	return fmt.Sprintf("💰 Платежи (%d)", count)
}
