package plans_transport_tg

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	decimal "github.com/shopspring/decimal"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
	"go.uber.org/zap"
)

func stepAwaitDurationUnit() core_tg_dialog.Step {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Месяцы", core_tg_dialog.ButtonData(durationUnitMonths)),
			tgbotapi.NewInlineKeyboardButtonData("Дни", core_tg_dialog.ButtonData(durationUnitDays)),
		),
	)

	return core_tg_dialog.Step{
		ID:     StepAwaitDurationUnit,
		Enter:  "➕ Новый срок\n\nВыберите единицу (или введите «м» / «д»):",
		Markup: &markup,
		OnInput: func(
			ctx context.Context,
			session *core_tg_dialog.Session,
			input string,
		) core_tg_dialog.Result {
			var message string
			switch strings.ToLower(strings.TrimSpace(input)) {
			case "месяцы", "м", durationUnitMonths:
				session.Data[dataKeyDurationUnit] = durationUnitMonths
				message = fmt.Sprintf("Введите количество месяцев (1–%d):", domain.PlanMaxMonths)
			case "дни", "д", durationUnitDays:
				session.Data[dataKeyDurationUnit] = durationUnitDays
				message = fmt.Sprintf("Введите количество дней (1–%d):", domain.PlanMaxDays)
			default:
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: "Нажмите кнопку или введите «м» / «д».",
				}
			}

			return core_tg_dialog.Result{
				Kind:    core_tg_dialog.Advance,
				Next:    StepAwaitDurationValue,
				Message: message,
			}
		},
	}
}

func stepAwaitDurationValue() core_tg_dialog.Step {
	return core_tg_dialog.Step{
		ID:    StepAwaitDurationValue,
		Enter: "Введите длительность:",
		OnInput: func(
			ctx context.Context,
			session *core_tg_dialog.Session,
			input string,
		) core_tg_dialog.Result {
			val, err := strconv.ParseInt(strings.TrimSpace(input), 10, 32)
			if err != nil {
				if errors.Is(err, strconv.ErrRange) {
					return planValidationResult(domain.ErrPlanDurationTooLong, "duration out of range")
				}
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: "Введите целое число.",
				}
			}

			var duration domain.Duration
			switch session.Data[dataKeyDurationUnit] {
			case durationUnitMonths:
				duration = domain.NewDuration(int32(val), 0, 0)
			case durationUnitDays:
				duration = domain.NewDuration(0, int32(val), 0)
			default:
				return core_tg_dialog.Result{
					Kind:   core_tg_dialog.Abort,
					Err:    fmt.Errorf("unknown duration unit in session: %q", session.Data[dataKeyDurationUnit]),
					LogMsg: "broken dialog session",
				}
			}

			if err := domain.ValidatePlanDuration(duration); err != nil {
				return planValidationResult(err, "invalid plan duration")
			}

			session.Data[dataKeyDurationMonths] = strconv.Itoa(int(duration.Months))
			session.Data[dataKeyDurationDays] = strconv.Itoa(int(duration.Days))

			return core_tg_dialog.Result{
				Kind:    core_tg_dialog.Advance,
				Next:    StepAwaitPrice,
				Message: fmt.Sprintf("Срок: %s.\nВведите цену в рублях (например 150 или 99.90):", formatPriceDuration(duration)),
			}
		},
	}
}

func stepAwaitPrice(plansService PlansService) core_tg_dialog.Step {
	return core_tg_dialog.Step{
		ID:    StepAwaitPrice,
		Enter: "Введите цену в рублях (например 150 или 99.90):",
		OnInput: func(
			ctx context.Context,
			session *core_tg_dialog.Session,
			input string,
		) core_tg_dialog.Result {
			text := strings.ReplaceAll(strings.TrimSpace(input), ",", ".")

			price, err := decimal.NewFromString(text)
			if err != nil {
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: "Введите число, например 150 или 99.90.",
				}
			}

			target, err := loadPlanTarget(session)
			if err != nil {
				return core_tg_dialog.Result{Kind: core_tg_dialog.Abort, Err: err, LogMsg: "broken dialog session"}
			}

			duration, err := loadDuration(session)
			if err != nil {
				return core_tg_dialog.Result{Kind: core_tg_dialog.Abort, Err: err, LogMsg: "broken dialog session"}
			}

			priceDomain, err := domain.NewPlanPrice(target.planID, duration, price)
			if err != nil {
				return planValidationResult(err, "invalid plan price")
			}

			_, err = plansService.CreatePlanPrice(ctx, priceDomain)
			switch {
			case errors.Is(err, core_errors.ErrConflict):
				core_logger.FromContext(ctx).Info("active plan price with same duration exists", zap.Error(err))
				return completeWithPlanCard(ctx, plansService, target,
					"⚠️ В тарифе уже есть активный срок такой же длительности "+
						"(1 месяц и 30 дней считаются одинаковыми).\n"+
						"Чтобы поменять цену — выключите старый срок и добавьте новый.",
				)
			case errors.Is(err, core_errors.ErrNotFound):
				markup := plansMenuMarkup()
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Complete,
					Message: "⚠️ Тариф не найден — возможно, его удалили.",
					Markup:  &markup,
				}
			case err != nil:
				return core_tg_dialog.Result{
					Kind:   core_tg_dialog.Abort,
					Err:    fmt.Errorf("create plan price: %w", err),
					LogMsg: "failed to create plan price",
				}
			}

			return completeWithPlanCard(ctx, plansService, target, "✅ Срок добавлен.")
		},
	}
}

// loadDuration — срок, сохранённый шагом stepAwaitDurationValue
func loadDuration(session *core_tg_dialog.Session) (domain.Duration, error) {
	months, err := strconv.Atoi(session.Data[dataKeyDurationMonths])
	if err != nil {
		return domain.Duration{}, fmt.Errorf("session duration_months: %w", err)
	}
	days, err := strconv.Atoi(session.Data[dataKeyDurationDays])
	if err != nil {
		return domain.Duration{}, fmt.Errorf("session duration_days: %w", err)
	}
	return domain.NewDuration(int32(months), int32(days), 0), nil
}

// completeWithPlanCard завершает диалог сообщением + свежей карточкой тарифа.
// Если карточку загрузить не удалось — сообщение остаётся, кнопка ведёт в список.
func completeWithPlanCard(
	ctx context.Context,
	plansService PlansService,
	target planTarget,
	header string,
) core_tg_dialog.Result {
	text, markup, err := loadPlanCard(ctx, plansService, target.planID, target.filter, target.page)
	if err != nil {
		core_logger.FromContext(ctx).Error("failed to load plan card after dialog", zap.Error(err))
		markup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(backToListButton(target.filter, target.page)),
		)
		return core_tg_dialog.Result{Kind: core_tg_dialog.Complete, Message: header, Markup: &markup}
	}

	return core_tg_dialog.Result{
		Kind:    core_tg_dialog.Complete,
		Message: header + "\n\n" + text,
		Markup:  &markup,
	}
}
