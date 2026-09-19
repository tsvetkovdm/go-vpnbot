package users_service

import (
	"context"
	"fmt"
)

func (s *UsersService) BlockUser(
	ctx context.Context,
	telegramID int64,
) error {
	if err := s.usersRepository.BlockUser(ctx, telegramID); err != nil {
		return fmt.Errorf("block user: %w", err)
	}

	return nil
}
