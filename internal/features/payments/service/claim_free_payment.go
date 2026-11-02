package payments_service

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
)

func (s *PaymentsService) ClaimFreePayment(
	ctx context.Context,
	telegramID int64,
	planPriceID int64,
) (domain.ResolvedPayment, error) {
	payment, err := s.preparePayment(ctx, telegramID, planPriceID)
	if err != nil {
		return domain.ResolvedPayment{}, fmt.Errorf("claim free payment: %w", err)
	}

	if !payment.IsFree() {
		return domain.ResolvedPayment{}, fmt.Errorf(
			"claim free payment: plan_price_id=%d is not free: %w",
			planPriceID, domain.ErrPaymentPlanUnavailable,
		)
	}

	payment.Status = domain.PaymentStatusSuccessful

	var result domain.ResolvedPayment
	err = s.transactor.WithinTx(ctx, func(ctx context.Context) error {
		created, err := s.paymentsRepository.CreatePayment(ctx, payment)
		if err != nil {
			return err
		}

		sub, err := s.subscriptionsService.Extend(ctx, created.UserID, created.Duration, 0)
		if errors.Is(err, core_errors.ErrConflict) {
			return fmt.Errorf("user_id=%d already had access: %w", created.UserID, domain.ErrPaymentFreeUnavailable)
		}
		if err != nil {
			return err
		}

		result = domain.ResolvedPayment{
			Payment:        created,
			UserTelegramID: telegramID,
			Subscription:   &sub,
		}
		return nil
	})
	if err != nil {
		return domain.ResolvedPayment{}, fmt.Errorf("claim free payment: %w", err)
	}

	return result, nil
}
