package core_tg_middleware

import (
	core_tg_handler "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/handler"
)

type Middleware func(core_tg_handler.Handler) core_tg_handler.Handler

func ChainMiddleware(
	h core_tg_handler.Handler,
	m ...Middleware,
) core_tg_handler.Handler {
	if len(m) == 0 {
		return h
	}

	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}
