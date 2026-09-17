package core_tg_response

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	return &TGResponseHandler{
		log:    log,
		update: update,
		bot:    bot,
	}
}

func (h *TGResponseHandler) TextResponse(text string) {
	chatID := core_tg_request.ChatID(h.update)

	msg := tgbotapi.NewMessage(chatID, text)

	if _, err := h.bot.Send(msg); err != nil {
		h.log.Error("write text response", zap.Error(err))
	}
}

func (h *TGResponseHandler) ErrorResponse(err error, msg string) {
	h.log.Error(msg, zap.Error(err))
	h.TextResponse("Что-то пошло не так, попробуйте позже.")
}

func (h *TGResponseHandler) PanicResponse(p any, msg string) {
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	h.TextResponse("Внутренняя ошибка сервиса. Пожалуйста попробуйте позже.")
}
