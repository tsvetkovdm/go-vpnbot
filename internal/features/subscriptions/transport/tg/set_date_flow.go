package subscriptions_transport_tg

import (
	"context"
	"errors"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
	"go.uber.org/zap"
)

const stepAwaitSetDateID core_tg_dialog.StepID = "set_date.await_date"

func newSetDateFlow(subscriptionsService SubscriptionsService) core_tg_dialog.Flow {
	exitMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ К пользователям", "users:menu"),
		),
	)

	return core_tg_dialog.NewFlow(
		"set_date",
		stepAwaitSetDateID,
		"Изменение даты отменено.",
		&exitMarkup,
		stepAwaitSetDate(subscriptionsService),
	)
}

func stepAwaitSetDate(subscriptionsService SubscriptionsService) core_tg_dialog.Step {
	return core_tg_dialog.Step{
		ID:    stepAwaitSetDateID,
		Enter: "📅 Введите, до какого дня оплачен доступ (ДД.ММ.ГГГГ или ДД.ММ):",
		OnInput: func(ctx context.Context, session *core_tg_dialog.Session, input string) core_tg_dialog.Result {
			date, err := parseDateInput(input, time.Now())
			if err != nil {
				return core_tg_dialog.Result{Kind: core_tg_dialog.Retry, Message: "Неверный формат даты. Пример: 12.11.2026 или 12.11"}
			}

			target, err := loadDialogTarget(session)
			if err != nil {
				return core_tg_dialog.Result{Kind: core_tg_dialog.Abort, Err: err, LogMsg: "broken dialog session"}
			}

			backToCard := backToCardMarkup(target.userID, target.filter, target.page)

			sub, err := subscriptionsService.SetPaidUntil(ctx, target.userID, date, target.version)
			switch {
			case errors.Is(err, domain.ErrSubscriptionDateInPast):
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: "Эта дата уже прошла.",
				}
			case errors.Is(err, domain.ErrSubscriptionDateTooFar):
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: fmt.Sprintf("Не дальше чем на %d лет вперёд.", domain.SubscriptionMaxYearsAhead),
				}
			case errors.Is(err, core_errors.ErrConflict):
				core_logger.FromContext(ctx).Info("subscription changed concurrently", zap.Error(err))
				return core_tg_dialog.Result{Kind: core_tg_dialog.Complete, Message: conflictText, Markup: &backToCard}
			case err != nil:
				return core_tg_dialog.Result{Kind: core_tg_dialog.Abort, Err: err, LogMsg: "failed to set paid until"}
			}

			return core_tg_dialog.Result{
				Kind: core_tg_dialog.Complete,
				Message: fmt.Sprintf(
					"✅ Доступ действует до %s.",
					sub.PaidUntil(time.Local).Format("02.01.2006"),
				),
				Markup: &backToCard,
			}
		},
	}
}
