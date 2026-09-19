package users_service

import (
	"context"
	"fmt"
)

func (s *UsersService) UnblockUser(
	ctx context.Context,
	telegramID int64,
) error {
	if err := s.usersRepository.UnblockUser(ctx, telegramID); err != nil {
		return fmt.Errorf("unblock user: %w", err)
	}

	return nil
}
