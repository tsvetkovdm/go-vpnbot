package core_tg_poller

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	commandRoutes  map[string]Route
	callbackRoutes []Route
}

func NewRouter() *Router {
	return &Router{
		commandRoutes: make(map[string]Route),
	}
}

func (r *Router) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		switch route.Kind {
		case TriggerCommand:
			r.commandRoutes[route.Trigger] = route
		case TriggerCallbackPrefix:
			r.callbackRoutes = append(r.callbackRoutes, route)
		}
	}
}

func (r *Router) Dispatch(ctx context.Context, update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		route, ok := r.commandRoutes[update.Message.Command()]
		if !ok {
			return
		}
		r.invoke(ctx, update, route)
	case update.CallbackQuery != nil:
		for i := range r.callbackRoutes {
			if strings.HasPrefix(update.CallbackQuery.Data, r.callbackRoutes[i].Trigger) {
				r.invoke(ctx, update, r.callbackRoutes[i])
				break
			}
		}
	}
}

func (r *Router) invoke(ctx context.Context, update tgbotapi.Update, route Route) {
	handler := route.WithMiddleware()
	handler(ctx, update)
}
