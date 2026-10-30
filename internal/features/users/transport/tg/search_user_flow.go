package users_transport_tg

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
)

const stepAwaitSearchQueryID core_tg_dialog.StepID = "search_user.await_query"

func newSearchUserFlow(usersService UsersService) core_tg_dialog.Flow {
	exitMarkup := backToUsersMenuMarkup()

	return core_tg_dialog.NewFlow(
		"search_user",
		stepAwaitSearchQueryID,
		"Поиск отменён.",
		&exitMarkup,
		stepAwaitSearchQuery(usersService),
	)
}

func stepAwaitSearchQuery(usersService UsersService) core_tg_dialog.Step {
	return core_tg_dialog.Step{
		ID:    stepAwaitSearchQueryID,
		Enter: "🔍 Введите @username или Telegram ID пользователя:",
		OnInput: func(
			ctx context.Context,
			session *core_tg_dialog.Session,
			input string,
		) core_tg_dialog.Result {
			query := strings.TrimPrefix(strings.TrimSpace(input), "@")
			if query == "" {
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: "Пустой запрос.",
				}
			}

			var overview domain.UserOverview
			var err error

			if telegramID, parseErr := strconv.ParseInt(query, 10, 64); parseErr == nil {
				overview, err = usersService.GetUserOverviewByTelegramID(ctx, telegramID)
			} else {
				overview, err = usersService.GetUserOverviewByUsername(ctx, query)
			}

			if errors.Is(err, core_errors.ErrNotFound) {
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: fmt.Sprintf("Пользователь «%s» не найден.", query),
				}
			}
			if err != nil {
				return core_tg_dialog.Result{
					Kind:   core_tg_dialog.Abort,
					Err:    err,
					LogMsg: "failed to search user",
				}
			}

			text, markup := assembleUserCard(overview, domain.UserFilterAll, 0, time.Now())

			return core_tg_dialog.Result{
				Kind:    core_tg_dialog.Complete,
				Message: text,
				Markup:  &markup,
			}
		},
	}
}
