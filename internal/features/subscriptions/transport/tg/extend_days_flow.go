package subscriptions_transport_tg

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
	"go.uber.org/zap"
)

const stepAwaitExtendDaysID core_tg_dialog.StepID = "extend_days.await_days"

func newExtendDaysFlow(subscriptionsService SubscriptionsService) core_tg_dialog.Flow {
	exitMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ К пользователям", "users:menu"),
		),
	)

	return core_tg_dialog.NewFlow(
		"extend_days",
		stepAwaitExtendDaysID,
		"Продление отменено.",
		&exitMarkup,
		stepAwaitExtendDays(subscriptionsService),
	)
}

func stepAwaitExtendDays(svc SubscriptionsService) core_tg_dialog.Step {
	return core_tg_dialog.Step{
		ID:    stepAwaitExtendDaysID,
		Enter: fmt.Sprintf("✏️ Введите, на сколько дней продлить (1–%d):", domain.SubscriptionMaxExtendDays),
		OnInput: func(ctx context.Context, session *core_tg_dialog.Session, input string) core_tg_dialog.Result {
			days, err := strconv.Atoi(strings.TrimSpace(input))
			if err != nil {
				return core_tg_dialog.Result{Kind: core_tg_dialog.Retry, Message: "Введите целое число."}
			}

			target, err := loadDialogTarget(session)
			if err != nil {
				return core_tg_dialog.Result{Kind: core_tg_dialog.Abort, Err: err, LogMsg: "broken dialog session"}
			}

			backToCard := backToCardMarkup(target.userID, target.filter, target.page)

			sub, err := svc.ExtendDays(ctx, target.userID, days, target.version)
			switch {
			case errors.Is(err, core_errors.ErrInvalidArgument):
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: fmt.Sprintf("Число дней должно быть от 1 до %d.", domain.SubscriptionMaxExtendDays),
				}
			case errors.Is(err, core_errors.ErrConflict):
				core_logger.FromContext(ctx).Info("subscription changed concurrently", zap.Error(err))
				return core_tg_dialog.Result{Kind: core_tg_dialog.Complete, Message: conflictText, Markup: &backToCard}
			case err != nil:
				return core_tg_dialog.Result{Kind: core_tg_dialog.Abort, Err: err, LogMsg: "failed to extend subscription"}
			}

			return core_tg_dialog.Result{
				Kind: core_tg_dialog.Complete,
				Message: fmt.Sprintf(
					"✅ Доступ продлён на %d дн.\nОплачено до %s",
					days,
					sub.PaidUntil(time.Local).Format("02.01.2006"),
				),
				Markup: &backToCard,
			}
		},
	}
}
