package payments_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PaymentsService) GetPaymentOverview(
	ctx context.Context,
	paymentID int64,
) (domain.PaymentOverview, error) {
	overview, err := s.paymentsRepository.GetPaymentOverview(ctx, paymentID)
	if err != nil {
		return domain.PaymentOverview{}, fmt.Errorf("get payment overview: %w", err)
	}

	return overview, nil
}
