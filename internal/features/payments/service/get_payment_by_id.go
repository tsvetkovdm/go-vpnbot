package payments_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PaymentsService) GetPaymentByID(
	ctx context.Context,
	paymentID int64,
) (domain.Payment, error) {
	payment, err := s.paymentsRepository.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return domain.Payment{}, fmt.Errorf("get payment by id: %w", err)
	}

	return payment, nil
}
