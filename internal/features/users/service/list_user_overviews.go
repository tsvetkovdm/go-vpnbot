package users_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *UsersService) ListUserOverviews(
	ctx context.Context,
	filter domain.UserFilter,
	limit, offset int,
) ([]domain.UserOverview, int, error) {
	overviews, total, err := s.usersRepository.ListUserOverviews(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list user overviews: %w", err)
	}

	return overviews, total, nil
}
