package access_transport_http

import (
	"context"
	"net/http"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

type AccessHTTPHandler struct {
	accessService AccessService
	config        Config
	adminContact  string
}

type AccessService interface {
	GetFeed(
		ctx context.Context,
		token string,
	) (domain.Feed, error)
}

func NewAccessHTTPHandler(
	accessService AccessService,
	config Config,
	adminContact string,
) *AccessHTTPHandler {
	return &AccessHTTPHandler{
		accessService: accessService,
		config:        config,
		adminContact:  adminContact,
	}
}

func (h *AccessHTTPHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /{token}", h.Feed)
}
