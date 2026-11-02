package payments_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *PaymentsService) ListPendingPayments(
	ctx context.Context,
	limit int,
) ([]domain.PaymentOverview, int, error) {
	overviews, total, err := s.paymentsRepository.ListPendingPayments(ctx, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("list pending payments: %w", err)
	}

	return overviews, total, nil
}
