package payments_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PaymentsService) ConfirmPayment(
	ctx context.Context,
	paymentID int64,
	adminTelegramID int64,
) (domain.ResolvedPayment, error) {
	var result domain.ResolvedPayment

	err := s.transactor.WithinTx(ctx, func(ctx context.Context) error {
		payment, userTelegramID, err := s.paymentsRepository.ResolvePayment(
			ctx, paymentID, adminTelegramID, domain.PaymentStatusSuccessful,
		)
		if err != nil {
			return err
		}

		sub, err := s.subscriptionsService.Extend(
			ctx, payment.UserID, payment.Duration, domain.SubscriptionSkipVersionCheck,
		)
		if err != nil {
			return err
		}

		result = domain.ResolvedPayment{
			Payment:        payment,
			UserTelegramID: userTelegramID,
			Subscription:   &sub,
		}
		return nil
	})
	if err != nil {
		return domain.ResolvedPayment{}, fmt.Errorf("confirm payment id=%d: %w", paymentID, err)
	}

	return result, nil
}
