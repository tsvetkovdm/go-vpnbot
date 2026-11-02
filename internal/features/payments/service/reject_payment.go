package payments_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PaymentsService) RejectPayment(
	ctx context.Context,
	paymentID int64,
	adminTelegramID int64,
) (domain.ResolvedPayment, error) {
	payment, userTelegramID, err := s.paymentsRepository.ResolvePayment(
		ctx, paymentID, adminTelegramID, domain.PaymentStatusRejected,
	)
	if err != nil {
		return domain.ResolvedPayment{}, fmt.Errorf("reject payment id=%d: %w", paymentID, err)
	}

	return domain.ResolvedPayment{
		Payment:        payment,
		UserTelegramID: userTelegramID,
	}, nil
}
