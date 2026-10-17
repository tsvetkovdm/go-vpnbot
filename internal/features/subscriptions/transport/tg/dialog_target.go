package subscriptions_transport_tg

import (
	"fmt"
	"strconv"

	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
)

const (
	dataKeyUserID  = "user_id"
	dataKeyVersion = "version"
	dataKeyFilter  = "filter"
	dataKeyPage    = "page"
)

type dialogTarget struct {
	userID  int64
	version int64
	filter  string
	page    string
}

// Из callback subs:<days|date>:<userID>:<version>:<filter>:<page>
func parseDialogCallback(data string) (dialogTarget, error) {
	segments := core_tg_callback.Segments(data)
	if len(segments) != 6 {
		return dialogTarget{}, fmt.Errorf("malformed callback data %q: %w", data, core_errors.ErrInvalidArgument)
	}

	userID, err := strconv.ParseInt(segments[2], 10, 64)
	if err != nil {
		return dialogTarget{}, fmt.Errorf("invalid userID %q: %w", segments[2], core_errors.ErrInvalidArgument)
	}

	version, err := strconv.ParseInt(segments[3], 10, 64)
	if err != nil {
		return dialogTarget{}, fmt.Errorf("invalid version %q: %w", segments[3], core_errors.ErrInvalidArgument)
	}

	return dialogTarget{
		userID:  userID,
		version: version,
		filter:  segments[4],
		page:    segments[5],
	}, nil
}

func (t dialogTarget) saveTo(session *core_tg_dialog.Session) {
	session.Data[dataKeyUserID] = strconv.FormatInt(t.userID, 10)
	session.Data[dataKeyVersion] = strconv.FormatInt(t.version, 10)
	session.Data[dataKeyFilter] = t.filter
	session.Data[dataKeyPage] = t.page
}

func loadDialogTarget(session *core_tg_dialog.Session) (dialogTarget, error) {
	userID, err := strconv.ParseInt(session.Data[dataKeyUserID], 10, 64)
	if err != nil {
		return dialogTarget{}, fmt.Errorf("session user_id: %w", err)
	}
	version, err := strconv.ParseInt(session.Data[dataKeyVersion], 10, 64)
	if err != nil {
		return dialogTarget{}, fmt.Errorf("session version: %w", err)
	}
	return dialogTarget{
		userID:  userID,
		version: version,
		filter:  session.Data[dataKeyFilter],
		page:    session.Data[dataKeyPage],
	}, nil
}
