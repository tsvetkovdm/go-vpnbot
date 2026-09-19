package users_service

import (
	"context"
	"fmt"
)

func (s *UsersService) DeleteUser(
	ctx context.Context,
	telegramID int64,
) error {
	if err := s.usersRepository.DeleteUser(ctx, telegramID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
