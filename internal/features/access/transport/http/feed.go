package access_transport_http

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	"go.uber.org/zap"
)

func (h *AccessHTTPHandler) Feed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	token := r.PathValue("token")
	if !isValidToken(token) {
		http.NotFound(w, r)
		return
	}

	feed, err := h.accessService.GetFeed(ctx, token)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		log.Error("get feed", zap.Error(err), zap.String("token", maskToken(token)))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	links := feed.Links()
	if len(links) == 0 {
		links = []string{placeholderLink(noAccessTitle(h.adminContact))}
	}

	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	header.Set("Cache-Control", "no-store")
	header.Set("Profile-Title", "base64:"+base64.StdEncoding.EncodeToString([]byte(h.config.ProfileTitle)))
	header.Set("Profile-Update-Interval", strconv.Itoa(h.config.UpdateIntervalHours))
	if feed.Owner.ExpiresAt != nil {
		header.Set("Subscription-Userinfo", fmt.Sprintf(
			"upload=0; download=0; total=0; expire=%d",
			feed.Owner.ExpiresAt.Unix()-1,
		))
	}
	if h.config.SupportURL != "" {
		header.Set("Support-Url", h.config.SupportURL)
	}

	body := base64.StdEncoding.EncodeToString([]byte(strings.Join(links, "\n")))
	if _, err := w.Write([]byte(body)); err != nil {
		log.Debug("write feed", zap.Error(err))
	}
}
