package core_tg_poller

import (
	core_tg_handler "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/handler"
	core_tg_middleware "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/middleware"
)

type TriggerKind string

const (
	TriggerCommand        TriggerKind = "command"
	TriggerCallbackPrefix TriggerKind = "callback_prefix"
)

type Route struct {
	Trigger    string
	Kind       TriggerKind
	Handler    core_tg_handler.Handler
	Middleware []core_tg_middleware.Middleware
}

func NewRoute(
	trigger string,
	kind TriggerKind,
	handler core_tg_handler.Handler,
	middleware ...core_tg_middleware.Middleware,
) Route {
	return Route{
		Trigger:    trigger,
		Kind:       kind,
		Handler:    handler,
		Middleware: middleware,
	}
}

func (r *Route) WithMiddleware() core_tg_handler.Handler {
	return core_tg_middleware.ChainMiddleware(
		r.Handler,
		r.Middleware...,
	)
}
