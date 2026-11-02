package core_tg_screen

import (
	"context"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	"go.uber.org/zap"
)

type Screens struct {
	bot  *tgbotapi.BotAPI
	mu   sync.Mutex
	data map[int64]int
}

func NewScreens(bot *tgbotapi.BotAPI) *Screens {
	return &Screens{
		bot:  bot,
		data: make(map[int64]int),
	}
}

func (s *Screens) Show(
	ctx context.Context,
	update tgbotapi.Update,
	text string,
	markup tgbotapi.InlineKeyboardMarkup,
) int {
	log := core_logger.FromContext(ctx)
	chatID := core_tg_request.ChatID(update)

	if update.Message != nil {
		s.deleteOrDisarm(log, chatID, update.Message.MessageID)
	}

	if oldID, ok := s.get(chatID); ok {
		if update.Message != nil && update.Message.MessageID == oldID+1 {
			edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, oldID, text, markup)
			_, err := s.bot.Send(edit)
			if err == nil || strings.Contains(err.Error(), "message is not modified") {
				return oldID
			}
		}
		s.deleteOrDisarm(log, chatID, oldID)
	}

	return s.send(log, chatID, text, markup)
}

func (s *Screens) Replace(
	ctx context.Context,
	chatID int64,
	text string,
	markup tgbotapi.InlineKeyboardMarkup,
) int {
	log := core_logger.FromContext(ctx)

	if oldID, ok := s.get(chatID); ok {
		s.deleteOrDisarm(log, chatID, oldID)
	}

	return s.send(log, chatID, text, markup)
}

func (s *Screens) send(
	log *core_logger.Logger,
	chatID int64,
	text string,
	markup tgbotapi.InlineKeyboardMarkup,
) int {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = markup

	sent, err := s.bot.Send(msg)
	if err != nil {
		log.Error("screen: send new screen", zap.Int64("chat_id", chatID), zap.Error(err))
		return 0
	}

	s.set(chatID, sent.MessageID)
	return sent.MessageID
}

func (s *Screens) deleteOrDisarm(log *core_logger.Logger, chatID int64, messageID int) {
	if _, err := s.bot.Request(tgbotapi.NewDeleteMessage(chatID, messageID)); err == nil {
		return
	}

	empty := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}}
	edit := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, empty)
	if _, err := s.bot.Request(edit); err != nil {
		log.Debug("screen: failed to delete or disarm message", zap.Int("message_id", messageID), zap.Error(err))
	}

}

func (s *Screens) get(chatID int64) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.data[chatID]
	return id, ok
}

func (s *Screens) set(chatID int64, messageID int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[chatID] = messageID
}
