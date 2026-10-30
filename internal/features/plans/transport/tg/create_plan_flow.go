package plans_transport_tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
)

const (
	StepAwaitTitle core_tg_dialog.StepID = "create_plan.await_title"
)

func newCreatePlanFlow(plansService PlansService) core_tg_dialog.Flow {
	exitMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("К тарифам", core_tg_callback.Build("plans", "menu")),
		),
	)

	return core_tg_dialog.NewFlow(
		"create_plan",
		StepAwaitTitle,
		"Создание тарифа отменено.",
		&exitMarkup,
		stepAwaitTitle(plansService),
	)
}
