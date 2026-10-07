package plans_transport_tg

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

const (
	plansPageSize = 6
)

func (h *PlansTGHandler) MenuListPlans(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuListPlans handler")

	data := update.CallbackQuery.Data
	segments := core_tg_callback.Segments(data)
	if len(segments) != 4 {
		responseHandler.ErrorResponse(
			fmt.Errorf("malformed callback data: %q", data),
			"malformed callback data",
		)
		return
	}

	includeInactive := strings.EqualFold(segments[2], "all")
	pageNumber, err := strconv.Atoi(segments[3])
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("parse page number %q: %w", segments[3], err),
			"failed to parse page number",
		)
		return
	}

	plans, err := h.plansService.ListPlans(ctx, includeInactive)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get list of plans")
		return
	}

	if len(plans) == 0 {
		data := core_tg_callback.Build("plans", "menu")
		markup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад", data),
			),
		)
		responseHandler.EditResponse("Тарифов не найдено.", &markup)
		return
	}

	maxPage := (len(plans) - 1) / plansPageSize
	if pageNumber > maxPage {
		pageNumber = maxPage
	}
	if pageNumber < 0 {
		pageNumber = 0
	}

	start := pageNumber * plansPageSize
	end := start + plansPageSize

	if end > len(plans) {
		end = len(plans)
	}

	pagePlans := plans[start:end]
	hasPrev := pageNumber > 0
	hasNext := end < len(plans)

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, plan := range pagePlans {
		text := fmt.Sprintf("%s - %v", plan.Title, plan.Price)
		callbackData := core_tg_callback.Build("plans", "view", strconv.FormatInt(plan.ID, 10), segments[2], strconv.Itoa(pageNumber))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(text, callbackData),
		))
	}

	navRow := []tgbotapi.InlineKeyboardButton{}

	if hasPrev {
		callbackData := core_tg_callback.Build("plans", "list", segments[2], strconv.Itoa(pageNumber-1))
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData("◀", callbackData))
	}

	if hasNext {
		callbackData := core_tg_callback.Build("plans", "list", segments[2], strconv.Itoa(pageNumber+1))
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData("▶", callbackData))
	}

	if len(navRow) > 0 {
		rows = append(rows, navRow)
	}

	callbackData := core_tg_callback.Build("plans", "list", "menu")
	backRow := []tgbotapi.InlineKeyboardButton{}
	backRow = append(backRow, tgbotapi.NewInlineKeyboardButtonData("Назад", callbackData))

	rows = append(rows, backRow)

	text := fmt.Sprintf("Страница %d из %d", pageNumber+1, maxPage+1)
	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)
	responseHandler.EditResponse(text, &markup)
}
