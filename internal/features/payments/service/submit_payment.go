package payments_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PaymentsService) SubmitPayment(
	ctx context.Context,
	telegramID int64,
	paymentID int64,
) (domain.Payment, error) {
	payment, err := s.paymentsRepository.SubmitPayment(ctx, telegramID, paymentID)
	if err != nil {
		return domain.Payment{}, fmt.Errorf("submit payment: %w", err)
	}

	return payment, nil
}
