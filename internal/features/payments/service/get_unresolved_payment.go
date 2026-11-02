package payments_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PaymentsService) GetUnresolvedPayment(
	ctx context.Context,
	telegramID int64,
) (domain.Payment, error) {
	payment, err := s.paymentsRepository.GetUnresolvedPaymentByTelegramID(ctx, telegramID)
	if err != nil {
		return domain.Payment{}, fmt.Errorf("get unresolved payment: %w", err)
	}

	return payment, nil
}
