package plans_transport_tg

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
)

func stepAwaitTitle(plansService PlansService) core_tg_dialog.Step {
	return core_tg_dialog.Step{
		ID:    StepAwaitTitle,
		Enter: fmt.Sprintf("Введите название тарифа (до %d символов):", domain.PlanTitleMaxLen),
		OnInput: func(
			ctx context.Context,
			session *core_tg_dialog.Session,
			input string,
		) core_tg_dialog.Result {
			planDomain, err := domain.NewPlan(input)
			if err != nil {
				return planValidationResult(err, "invalid plan title")
			}

			plan, err := plansService.CreatePlan(ctx, planDomain)
			if err != nil {
				return core_tg_dialog.Result{
					Kind:   core_tg_dialog.Abort,
					Err:    fmt.Errorf("create plan: %w", err),
					LogMsg: "failed to create plan",
				}
			}

			text, markup := assemblePlanCard(plan, nil, "all", "0")

			return core_tg_dialog.Result{
				Kind:    core_tg_dialog.Complete,
				Message: "✅ Тариф создан.\n\n" + text,
				Markup:  &markup,
			}
		},
	}
}
