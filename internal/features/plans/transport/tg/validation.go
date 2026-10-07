package plans_transport_tg

import (
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
)

func planValidationMessage(err error) (string, bool) {
	switch {
	case errors.Is(err, domain.ErrPlanTitleEmpty):
		return "Название не может быть пустым.", true
	case errors.Is(err, domain.ErrPlanTitleTooLong):
		return fmt.Sprintf("Превышен лимит символов. Максимальная длина: %d.", domain.PlanTitleMaxLen), true
	case errors.Is(err, domain.ErrPlanDurationInvalid):
		return "Длительность должна быть больше нуля.", true
	case errors.Is(err, domain.ErrPlanDurationTooLong):
		return fmt.Sprintf("Максимальная длительность: %d мес. или %d дн.", domain.PlanMaxMonths, domain.PlanMaxDays), true
	case errors.Is(err, domain.ErrPlanPriceNegative):
		return "Стоимость не может быть отрицательной.", true
	case errors.Is(err, domain.ErrPlanPriceTooHigh):
		return fmt.Sprintf("Стоимость не должна превышать %s р.", domain.PlanMaxPrice), true
	case errors.Is(err, domain.ErrPlanPricePrecision):
		return "Не более двух знаков после запятой.", true
	}
	return "", false
}

func planValidationResult(err error, logMsg string) core_tg_dialog.Result {
	if msg, ok := planValidationMessage(err); ok {
		return core_tg_dialog.Result{Kind: core_tg_dialog.Retry, Message: msg}
	}
	return core_tg_dialog.Result{Kind: core_tg_dialog.Abort, Err: err, LogMsg: logMsg}
}
