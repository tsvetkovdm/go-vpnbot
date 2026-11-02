package payments_transport_tg

import (
	"fmt"
	"strconv"

	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
)

func parseIDArg(data string) (int64, error) {
	segments := core_tg_callback.Segments(data)
	if len(segments) != 3 {
		return 0, fmt.Errorf("malformed callback data: %q", data)
	}

	id, err := strconv.ParseInt(segments[2], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse ID %q: %w", segments[2], err)
	}

	return id, nil
}

func callbackAction(data string) string {
	segments := core_tg_callback.Segments(data)
	if len(segments) < 2 {
		return ""
	}
	return segments[1]
}

func paymentCallback(action string, paymentID int64) string {
	return core_tg_callback.Build("pay", action, strconv.FormatInt(paymentID, 10))
}
