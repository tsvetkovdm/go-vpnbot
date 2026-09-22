package plans_transport_tg

import core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"

const (
	StepAwaitTitle         core_tg_dialog.StepID = "create_plan.await_title"
	StepAwaitDurationUnit  core_tg_dialog.StepID = "create_plan.await_duration_unit"
	StepAwaitDurationValue core_tg_dialog.StepID = "create_plan.await_duration_value"
	StepAwaitPrice         core_tg_dialog.StepID = "create_plan.await_price"
)

func newCreatePlanFlow(plansService PlansService) core_tg_dialog.Flow {
	return core_tg_dialog.NewFlow(
		"create_plan",
		StepAwaitTitle,
		stepAwaitTitle(),
		stepAwaitDurationUnit(),
		stepAwaitDurationValue(),
		stepAwaitPrice(plansService),
	)
}
