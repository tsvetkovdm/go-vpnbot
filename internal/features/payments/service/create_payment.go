package payments_service

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
)

func (s *PaymentsService) CreatePayment(
	ctx context.Context,
	telegramID int64,
	planPriceID int64,
) (domain.Payment, error) {
	payment, err := s.preparePayment(ctx, telegramID, planPriceID)
	if err != nil {
		return domain.Payment{}, fmt.Errorf("create payment: %w", err)
	}

	if payment.IsFree() {
		return domain.Payment{}, fmt.Errorf(
			"create payment: plan_price_id=%d is free: %w",
			planPriceID, domain.ErrPaymentPlanUnavailable,
		)
	}

	payment, err = s.paymentsRepository.CreatePayment(ctx, payment)
	if err != nil {
		return domain.Payment{}, fmt.Errorf("create payment: %w", err)
	}

	return payment, nil
}

func (s *PaymentsService) preparePayment(
	ctx context.Context,
	telegramID int64,
	planPriceID int64,
) (domain.Payment, error) {
	userID, status, err := s.paymentsRepository.GetUserStatusByTelegramID(ctx, telegramID)
	if err != nil {
		return domain.Payment{}, err
	}

	if status == domain.StatusBlocked {
		return domain.Payment{}, fmt.Errorf("user_id=%d: %w", userID, domain.ErrPaymentUserBlocked)
	}

	plan, price, err := s.paymentsRepository.GetPlanPriceForPayment(ctx, planPriceID)
	if errors.Is(err, core_errors.ErrNotFound) {
		return domain.Payment{}, fmt.Errorf(
			"plan_price_id=%d: %w", planPriceID, domain.ErrPaymentPlanUnavailable,
		)
	}
	if err != nil {
		return domain.Payment{}, err
	}

	return domain.NewPayment(userID, plan, price)
}
