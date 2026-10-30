package users_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *UsersService) GetUserOverviewByUsername(
	ctx context.Context,
	username string,
) (domain.UserOverview, error) {
	overview, err := s.usersRepository.GetUserOverviewByUsername(ctx, username)
	if err != nil {
		return domain.UserOverview{}, fmt.Errorf("get user overview by username: %w", err)
	}

	return overview, nil
}
