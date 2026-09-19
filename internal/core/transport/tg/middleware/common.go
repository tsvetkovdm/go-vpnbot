package core_tg_middleware

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_handler "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/handler"
	core_tg_request "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/request"
	core_tg_response "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/response"
	"go.uber.org/zap"
)

func Logger(log *core_logger.Logger) Middleware {
	return func(next core_tg_handler.Handler) core_tg_handler.Handler {
		return func(ctx context.Context, update tgbotapi.Update) {
			chatID := core_tg_request.ChatID(update)

			l := log.With(
				zap.Int64("chat_id", chatID),
				zap.Int("update_id", update.UpdateID),
			)

			ctx = core_logger.ToContext(ctx, l)

			next(ctx, update)
		}
	}
}

func Trace() Middleware {
	return func(next core_tg_handler.Handler) core_tg_handler.Handler {
		return func(ctx context.Context, update tgbotapi.Update) {
			log := core_logger.FromContext(ctx)

			before := time.Now()

			var command, args string
			if update.Message != nil {
				command = update.Message.Command()
				args = update.Message.CommandArguments()
			} else if update.CallbackQuery != nil {
				command = "callback"
				args = update.CallbackQuery.Data
			}

			log.Debug(
				">>> incoming tg request",
				zap.String("command", command),
				zap.String("arguments", args),
				zap.Time("time", before.UTC()),
			)

			next(ctx, update)

			log.Debug(
				"<<< done tg request",
				zap.Duration("latency", time.Since(before)),
			)
		}
	}
}

func Panic(bot *tgbotapi.BotAPI) Middleware {
	return func(next core_tg_handler.Handler) core_tg_handler.Handler {
		return func(ctx context.Context, update tgbotapi.Update) {
			defer func() {
				if p := recover(); p != nil {
					log := core_logger.FromContext(ctx)
					responseHandler := core_tg_response.NewTGResponseHandler(log, update, bot)
					responseHandler.PanicResponse(
						p,
						"during handle tg request go unexpected panic",
					)
				}
			}()

			next(ctx, update)
		}
	}
}
