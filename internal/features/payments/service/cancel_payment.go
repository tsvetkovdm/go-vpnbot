package payments_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PaymentsService) CancelPayment(
	ctx context.Context,
	telegramID int64,
	paymentID int64,
) (domain.Payment, error) {
	payment, err := s.paymentsRepository.CancelPayment(ctx, telegramID, paymentID)
	if err != nil {
		return domain.Payment{}, fmt.Errorf("cancel payment: %w", err)
	}

	return payment, nil
}
