package plans_transport_tg

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shopspring/decimal"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
)

const (
	durationUnitMonths string = "months"
	durationUnitDays   string = "days"
)

func stepAwaitTitle() core_tg_dialog.Step {
	return core_tg_dialog.Step{
		ID:    StepAwaitTitle,
		Enter: fmt.Sprintf("Введите название тарифа (до %d символов):", domain.PlanTitleMaxLen),
		OnInput: func(
			ctx context.Context,
			session *core_tg_dialog.Session,
			input string,
		) core_tg_dialog.Result {
			text := strings.TrimSpace(input)

			if err := domain.ValidatePlanTitle(text); err != nil {
				return planValidationResult(err, "invalid plan title")
			}

			session.Data[dataKeyTitle] = text

			return core_tg_dialog.Result{
				Kind: core_tg_dialog.Advance,
				Next: StepAwaitDurationUnit,
			}
		},
	}
}

func stepAwaitDurationUnit() core_tg_dialog.Step {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Месяцы", core_tg_dialog.ButtonData(durationUnitMonths)),
			tgbotapi.NewInlineKeyboardButtonData("Дни", core_tg_dialog.ButtonData(durationUnitDays)),
		),
	)

	return core_tg_dialog.Step{
		ID:     StepAwaitDurationUnit,
		Enter:  "Выберите единицу длительности (или введите «м» / «д»):",
		Markup: &markup,
		OnInput: func(
			ctx context.Context,
			session *core_tg_dialog.Session,
			input string,
		) core_tg_dialog.Result {
			text := strings.Fields(input)

			if len(text) != 1 {
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: "Некорректный ввод. Введите первую букву или слово целиком (Пример: месяцы).",
				}
			}

			var message string
			switch strings.ToLower(text[0]) {
			case "месяцы", "м", durationUnitMonths:
				session.Data[dataKeyDurationUnit] = durationUnitMonths
				message = "Введите длительность тарифа в месяцах:"
			case "дни", "д", durationUnitDays:
				session.Data[dataKeyDurationUnit] = durationUnitDays
				message = "Введите длительность тарифа в днях:"
			default:
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: "Некорректный ввод. Введите первую букву или слово целиком (Пример: месяцы):",
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
		Enter: "Введите длительность тарифа в месяцах/днях:",
		OnInput: func(
			ctx context.Context,
			session *core_tg_dialog.Session,
			input string,
		) core_tg_dialog.Result {
			text := strings.TrimSpace(input)

			val, err := strconv.ParseInt(text, 10, 32)
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
					LogMsg: "failed to create plan",
				}
			}

			if err := domain.ValidatePlanDuration(duration); err != nil {
				return planValidationResult(err, "invalid plan duration")
			}

			session.Data[dataKeyDurationMonths] = strconv.Itoa(int(duration.Months))
			session.Data[dataKeyDurationDays] = strconv.Itoa(int(duration.Days))

			return core_tg_dialog.Result{
				Kind: core_tg_dialog.Advance,
				Next: StepAwaitPrice,
			}
		},
	}
}

func stepAwaitPrice(plansService PlansService) core_tg_dialog.Step {
	return core_tg_dialog.Step{
		ID:    StepAwaitPrice,
		Enter: "Введите стоимость тарифа в рублях:",
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

			months, err := strconv.Atoi(session.Data[dataKeyDurationMonths])
			if err != nil {
				return core_tg_dialog.Result{
					Kind:   core_tg_dialog.Abort,
					Err:    fmt.Errorf("atoi months: %w", err),
					LogMsg: "failed to create plan",
				}
			}
			days, err := strconv.Atoi(session.Data[dataKeyDurationDays])
			if err != nil {
				return core_tg_dialog.Result{
					Kind:   core_tg_dialog.Abort,
					Err:    fmt.Errorf("atoi days: %w", err),
					LogMsg: "failed to create plan",
				}
			}

			planDomain, err := domain.NewPlan(
				session.Data[dataKeyTitle],
				domain.NewDuration(int32(months), int32(days), 0),
				price,
			)
			if err != nil {
				return planValidationResult(err, "failed to build plan")
			}

			plan, err := plansService.CreatePlan(ctx, planDomain)
			if err != nil {
				return core_tg_dialog.Result{
					Kind:   core_tg_dialog.Abort,
					Err:    fmt.Errorf("create plan: %w", err),
					LogMsg: "failed to create plan",
				}
			}

			duration := formatDuration(plan)

			message := fmt.Sprintf(
				"Тариф '%s' (%v) создан успешно. Цена: %v р.",
				plan.Title,
				duration,
				plan.Price,
			)

			return core_tg_dialog.Result{
				Kind:    core_tg_dialog.Complete,
				Message: message,
			}
		},
	}
}
