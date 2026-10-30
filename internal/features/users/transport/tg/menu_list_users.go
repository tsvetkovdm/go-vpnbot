package users_transport_tg

import (
	"context"
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
)

func parseListCallback(data string) (domain.UserFilter, int, error) {
	segments := core_tg_callback.Segments(data)
	if len(segments) != 4 {
		return "", 0, fmt.Errorf("malformed callback data %q: %w", data, core_errors.ErrInvalidArgument)
	}

	filter := domain.UserFilter(segments[2])
	if _, ok := filterTitles[filter]; !ok {
		return "", 0, fmt.Errorf("unknown filter %q: %w", segments[2], core_errors.ErrInvalidArgument)
	}

	page, err := strconv.Atoi(segments[3])
	if err != nil || page < 0 {
		return "", 0, fmt.Errorf("invalid page %q: %w", segments[3], core_errors.ErrInvalidArgument)
	}

	return filter, page, nil
}

func (h *UsersTGHandler) MenuListUsers(
	ctx context.Context,
	update tgbotapi.Update,
) {
	log := core_logger.FromContext(ctx)
	responseHandler := core_tg_response.NewTGResponseHandler(log, update, h.bot)

	log.Debug("invoke MenuListUsers handler")

	messageID := update.CallbackQuery.Message.MessageID
	backMarkup := backToUsersMenuMarkup()

	filter, page, err := parseListCallback(update.CallbackQuery.Data)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to parse callback data", &backMarkup)
		return
	}

	overviews, total, err := h.usersService.ListUserOverviews(ctx, filter, usersPageSize, page*usersPageSize)
	if err != nil {
		responseHandler.EditErrorResponse(messageID, err, "failed to list users", &backMarkup)
		return
	}

	if len(overviews) == 0 && page > 0 {
		page = 0
		overviews, total, err = h.usersService.ListUserOverviews(ctx, filter, usersPageSize, 0)
		if err != nil {
			responseHandler.EditErrorResponse(messageID, err, "failed to list users", &backMarkup)
			return
		}
	}

	title := filterTitles[filter]

	if total == 0 {
		responseHandler.EditResponse(title+"\n\nНикого нет.", &backMarkup)
		return
	}

	now := time.Now()
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, o := range overviews {
		data := core_tg_callback.Build(
			"users", "view",
			strconv.FormatInt(o.User.ID, 10),
			string(filter),
			strconv.Itoa(page),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(userButtonLabel(o, now), data),
		))
	}

	maxPage := (total - 1) / usersPageSize

	var nav []tgbotapi.InlineKeyboardButton
	if page > 0 {
		nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("◀️", listCallback(filter, page-1)))
	}
	if page < maxPage {
		nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("▶️", listCallback(filter, page+1)))
	}
	if len(nav) > 0 {
		rows = append(rows, nav)
	}

	text := fmt.Sprintf("%s - %d\nСтраница %d из %d", title, total, page+1, maxPage+1)
	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)
	markup.InlineKeyboard = append(markup.InlineKeyboard, backMarkup.InlineKeyboard...)
	responseHandler.EditResponse(text, &markup)
}
