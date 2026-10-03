package core_tg_response

import (
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	"go.uber.org/zap"
)

type TGResponseHandler struct {
	log    *core_logger.Logger
	update tgbotapi.Update
	bot    *tgbotapi.BotAPI
}

func NewTGResponseHandler(
	log *core_logger.Logger,
	update tgbotapi.Update,
	bot *tgbotapi.BotAPI,
) *TGResponseHandler {
	h := &TGResponseHandler{
		log:    log,
		update: update,
		bot:    bot,
	}

	if update.CallbackQuery != nil {
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		if _, err := bot.Request(callback); err != nil {
			log.Error("answer callback query", zap.Error(err))
		}
	}

	return h
}

func (h *TGResponseHandler) TextResponse(text string) int {
	chatID := core_tg_request.ChatID(h.update)

	messageConfig := tgbotapi.NewMessage(chatID, text)

	message, err := h.bot.Send(messageConfig)
	if err != nil {
		h.log.Error("write text response", zap.Error(err))
		return 0
	}

	return message.MessageID
}

func (h *TGResponseHandler) KeyboardResponse(
	text string,
	markup tgbotapi.InlineKeyboardMarkup,
) int {
	chatID := core_tg_request.ChatID(h.update)

	messageConfig := tgbotapi.NewMessage(chatID, text)
	messageConfig.ReplyMarkup = markup

	message, err := h.bot.Send(messageConfig)
	if err != nil {
		h.log.Error("write keyboard response", zap.Error(err))
		return 0
	}

	return message.MessageID
}

func (h *TGResponseHandler) EditResponse(
	text string,
	markup *tgbotapi.InlineKeyboardMarkup,
) {
	if h.update.CallbackQuery == nil || h.update.CallbackQuery.Message == nil {
		h.log.Error("edit response: update has no callback query or message")
		return
	}

	chatID := core_tg_request.ChatID(h.update)
	messageID := h.update.CallbackQuery.Message.MessageID

	h.EditMessage(chatID, messageID, text, markup)

}

func (h *TGResponseHandler) EditMessage(
	chatID int64,
	messageID int,
	text string,
	markup *tgbotapi.InlineKeyboardMarkup,
) {
	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	if markup != nil {
		edit.ReplyMarkup = markup
	}

	if _, err := h.bot.Send(edit); err != nil {
		var apiErr *tgbotapi.Error
		if errors.As(err, &apiErr) && strings.Contains(apiErr.Message, "message is not modified") {
			h.log.Debug(apiErr.Message, zap.Error(err))
			return
		}
		h.log.Error("edit response", zap.Error(err))
	}
}

func (h *TGResponseHandler) resolveError(err error) (string, func(string, ...zap.Field)) {
	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		return "Некорректные данные, попробуйте ещё раз.", h.log.Warn
	case errors.Is(err, core_errors.ErrNotFound):
		return "Запись не найдена.", h.log.Debug
	default:
		return "Что-то пошло не так, попробуйте позже.", h.log.Error
	}
}

func (h *TGResponseHandler) ErrorResponse(err error, msg string) {
	userText, logFunc := h.resolveError(err)
	logFunc(msg, zap.Error(err))
	h.TextResponse(userText)
}

func (h *TGResponseHandler) EditErrorResponse(chatID int64, messageID int, err error, msg string) {
	userText, logFunc := h.resolveError(err)
	logFunc(msg, zap.Error(err))
	h.EditMessage(chatID, messageID, userText, nil)
}

func (h *TGResponseHandler) PanicResponse(p any, msg string) {
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	h.TextResponse("Внутренняя ошибка сервиса. Пожалуйста попробуйте позже.")
}
