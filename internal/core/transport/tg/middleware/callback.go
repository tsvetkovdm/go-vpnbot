package core_tg_middleware

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_handler "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/handler"
	"go.uber.org/zap"
)

func CallbackAnswer(bot *tgbotapi.BotAPI) Middleware {
	return func(next core_tg_handler.Handler) core_tg_handler.Handler {
		return func(ctx context.Context, update tgbotapi.Update) {
			if update.CallbackQuery != nil {
				var answerText func() string
				ctx, answerText = core_tg_callback.WithAnswer(ctx)

				defer func() {
					callback := tgbotapi.NewCallback(update.CallbackQuery.ID, answerText())
					if _, err := bot.Request(callback); err != nil {
						core_logger.FromContext(ctx).Debug("answer callback query", zap.Error(err))
					}
				}()
			}

			next(ctx, update)
		}
	}
}
