package plans_transport_tg

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shopspring/decimal"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
)

const (
	durationUnitMonths string = "months"
	durationUnitDays   string = "days"

	maxPriceRubles = 1_000_000
)

func stepAwaitTitle() core_tg_dialog.Step {
	return core_tg_dialog.Step{
		ID:    StepAwaitTitle,
		Enter: "Введите название тарифа (до 30 символов)",
		OnInput: func(
			ctx context.Context,
			session *core_tg_dialog.Session,
			update tgbotapi.Update,
		) core_tg_dialog.Result {
			text := strings.TrimSpace(update.Message.Text)

			if utf8.RuneCountInString(text) > 30 {
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: "Превышен лимит символов. Максимальная длина: 30 символов",
				}
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
	return core_tg_dialog.Step{
		ID:    StepAwaitDurationUnit,
		Enter: "Выберите 'месяцы' или 'дни' (доступен текстовый ввод)",
		OnInput: func(
			ctx context.Context,
			session *core_tg_dialog.Session,
			update tgbotapi.Update,
		) core_tg_dialog.Result {
			text := strings.Fields(update.Message.Text)

			if len(text) != 1 {
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: "Некорректный ввод. Введите первую букву или слово целиком (Пример: месяцы).",
				}
			}

			var message string
			switch strings.ToLower(text[0]) {
			case "месяцы", "м":
				session.Data[dataKeyDurationUnit] = durationUnitMonths
				message = "Введите длительность тарифа в месяцах"
			case "дни", "д":
				session.Data[dataKeyDurationUnit] = durationUnitDays
				message = "Введите длительность тарифа в днях"
			default:
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: "Некорректный ввод. Введите первую букву или слово целиком (Пример: месяцы).",
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
		Enter: "Введите длительность тарифа в месяцах/днях",
		OnInput: func(
			ctx context.Context,
			session *core_tg_dialog.Session,
			update tgbotapi.Update,
		) core_tg_dialog.Result {
			text := strings.TrimSpace(update.Message.Text)

			val, err := strconv.Atoi(text)
			if err != nil {
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: "Введенное значение не является целым числом.",
				}
			}

			if val <= 0 {
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: "Количество должно быть больше нуля.",
				}
			}

			switch session.Data[dataKeyDurationUnit] {
			case durationUnitMonths:
				if val > 24 {
					return core_tg_dialog.Result{
						Kind:    core_tg_dialog.Retry,
						Message: "Количество месяцев не должно превышать 24.",
					}
				}
				session.Data[dataKeyDurationMonths] = text
				session.Data[dataKeyDurationDays] = "0"
			case durationUnitDays:
				if val > 366 {
					return core_tg_dialog.Result{
						Kind:    core_tg_dialog.Retry,
						Message: "Количество дней не должно превышать 366.",
					}
				}
				session.Data[dataKeyDurationDays] = text
				session.Data[dataKeyDurationMonths] = "0"
			}

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
		Enter: "Введите стоимость тарифа в рублях",
		OnInput: func(
			ctx context.Context,
			session *core_tg_dialog.Session,
			update tgbotapi.Update,
		) core_tg_dialog.Result {
			text := strings.TrimSpace(update.Message.Text)
			text = strings.ReplaceAll(text, ",", ".")

			price, err := decimal.NewFromString(text)
			if err != nil {
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: "Введенное значение не является целым/дробным числом.",
				}
			}

			if price.LessThan(decimal.NewFromInt(0)) {
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: "Стоимость не может быть отрицательной.",
				}
			}

			if price.GreaterThan(decimal.NewFromInt(maxPriceRubles)) {
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: fmt.Sprintf("Стоимость не должна превышать %d р.", maxPriceRubles),
				}
			}

			if price.Exponent() < -2 {
				return core_tg_dialog.Result{
					Kind:    core_tg_dialog.Retry,
					Message: "Не более двух знаков после запятой.",
				}
			}

			title := session.Data[dataKeyTitle]
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

			duration := domain.NewDuration(int32(months), int32(days), 0)

			planDomain := domain.NewPlan(title, price, duration)

			plan, err := plansService.CreatePlan(ctx, planDomain)
			if err != nil {
				return core_tg_dialog.Result{
					Kind:   core_tg_dialog.Abort,
					Err:    fmt.Errorf("create plan: %w", err),
					LogMsg: "failed to create plan",
				}
			}

			var durationStr string
			if plan.Duration.Months > 0 {
				durationStr = fmt.Sprintf("%d мес.", plan.Duration.Months)
			} else {
				durationStr = fmt.Sprintf("%d дн.", plan.Duration.Days)
			}

			message := fmt.Sprintf(
				"Тариф '%s' (%v) создан успешно. Цена: %v р.",
				plan.Title,
				durationStr,
				plan.Price,
			)

			return core_tg_dialog.Result{
				Kind:    core_tg_dialog.Complete,
				Message: message,
			}
		},
	}
}
