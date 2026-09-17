package core_tg_poller

import (
	"context"
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_middleware "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/middleware"
)

type TGPoller struct {
	router     *Router
	bot        *tgbotapi.BotAPI
	config     Config
	log        *core_logger.Logger
	middleware []core_tg_middleware.Middleware
}

func NewTGPoller(
	router *Router,
	bot *tgbotapi.BotAPI,
	config Config,
	log *core_logger.Logger,
	middleware ...core_tg_middleware.Middleware,
) *TGPoller {
	return &TGPoller{
		router:     router,
		bot:        bot,
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (p *TGPoller) Run(ctx context.Context) error {
	var wg = sync.WaitGroup{}
	handler := core_tg_middleware.ChainMiddleware(p.router.Dispatch, p.middleware...)

	updateConfig := tgbotapi.NewUpdate(-1)
	updateConfig.Timeout = int(p.config.PollTimeout)

	p.log.Warn("start TG poller")
	ch := p.bot.GetUpdatesChan(updateConfig)

	for {
		select {
		case update := <-ch:
			wg.Add(1)
			go func(ctx context.Context, update tgbotapi.Update) {
				defer wg.Done()

				handler(ctx, update)

			}(ctx, update)

		case <-ctx.Done():
			p.log.Warn("shutdown TG poller...")
			p.bot.StopReceivingUpdates()

			shutdownCtx, cancel := context.WithTimeout(
				context.Background(),
				p.config.ShutdownTimeout,
			)
			defer cancel()

			doneChan := make(chan struct{})

			go func() {
				wg.Wait()
				close(doneChan)
			}()

			select {
			case <-doneChan:
				p.log.Warn("all handlers finished gracefully")
			case <-shutdownCtx.Done():
				return fmt.Errorf("shutdown TG poller: %w", shutdownCtx.Err())
			}

			p.log.Warn("TG poller stopped")

			return nil
		}
	}
}
