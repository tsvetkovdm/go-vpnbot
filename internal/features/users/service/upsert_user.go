package users_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *UsersService) UpsertUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	user, err := s.usersRepository.UpsertUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("upsert user: %w", err)
	}

	return user, nil
}
