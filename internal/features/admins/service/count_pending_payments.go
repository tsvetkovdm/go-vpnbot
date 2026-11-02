package admins_service

import (
	"context"
	"fmt"
)

func (s *AdminsService) CountPendingPayments(
	ctx context.Context,
) (int, error) {
	count, err := s.adminsRepository.CountPendingPayments(ctx)
	if err != nil {
		return 0, fmt.Errorf("count pending payments: %w", err)
	}

	return count, nil
}
