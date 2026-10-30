package users_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *UsersService) GetUserOverview(
	ctx context.Context,
	userID int64,
) (domain.UserOverview, error) {
	overview, err := s.usersRepository.GetUserOverview(ctx, userID)
	if err != nil {
		return domain.UserOverview{}, fmt.Errorf("get user overview: %w", err)
	}

	return overview, nil
}
