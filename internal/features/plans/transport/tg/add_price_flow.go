package plans_transport_tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
)

const (
	StepAwaitDurationUnit  core_tg_dialog.StepID = "add_price.await_duration_unit"
	StepAwaitDurationValue core_tg_dialog.StepID = "add_price.await_duration_value"
	StepAwaitPrice         core_tg_dialog.StepID = "add_price.await_price"
)

const (
	dataKeyDurationUnit   = "duration_unit"
	dataKeyDurationMonths = "duration_months"
	dataKeyDurationDays   = "duration_days"
)

const (
	durationUnitMonths = "months"
	durationUnitDays   = "days"
)

func newAddPriceFlow(plansService PlansService) core_tg_dialog.Flow {
	exitMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ К тарифам", core_tg_callback.Build("plans", "menu")),
		),
	)

	return core_tg_dialog.NewFlow(
		"add_price",
		StepAwaitDurationUnit,
		"Добавление срока отменено.",
		&exitMarkup,
		stepAwaitDurationUnit(),
		stepAwaitDurationValue(),
		stepAwaitPrice(plansService),
	)
}
