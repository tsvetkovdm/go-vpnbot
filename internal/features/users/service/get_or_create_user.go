package users_service

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (s *UsersService) GetOrCreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	created, err := s.usersRepository.CreateUser(ctx, user)
	if err == nil {
		return created, nil
	}
	if errors.Is(err, core_postgres_pool.ErrNoRows) {
		return s.usersRepository.GetUserByTelegramID(ctx, user.TelegramID)
	}
	return domain.User{}, fmt.Errorf("get or create user: %w", err)
}
