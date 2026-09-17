package users_service

import (
	"context"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

type UsersService struct {
	usersRepository UsersRepository
}

type UsersRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUserByTelegramID(
		ctx context.Context,
		telegramID int64,
	) (domain.User, error)
}

func NewUsersService(
	usersRepository UsersRepository,
) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
