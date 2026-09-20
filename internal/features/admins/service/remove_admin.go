package admins_service

import (
	"context"
	"fmt"
)

func (s *AdminsService) RemoveAdmin(
	ctx context.Context,
	telegramID int64,
) error {
	if err := s.adminsRepository.RemoveAdmin(ctx, telegramID); err != nil {
		return fmt.Errorf("remove admin: %w", err)
	}

	return nil
}
