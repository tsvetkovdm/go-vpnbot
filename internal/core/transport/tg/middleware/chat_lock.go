package core_tg_middleware

import (
	"context"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_tg_handler "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/handler"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
)

func ChatLock() Middleware {
	var locks sync.Map

	return func(next core_tg_handler.Handler) core_tg_handler.Handler {
		return func(ctx context.Context, update tgbotapi.Update) {
			chatID := core_tg_request.ChatID(update)

			value, _ := locks.LoadOrStore(chatID, &sync.Mutex{})
			mu := value.(*sync.Mutex)

			mu.Lock()
			defer mu.Unlock()

			next(ctx, update)
		}
	}
}
