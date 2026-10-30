package users_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *UsersService) GetUserOverviewByTelegramID(
	ctx context.Context,
	telegramID int64,
) (domain.UserOverview, error) {
	overview, err := s.usersRepository.GetUserOverviewByTelegramID(ctx, telegramID)
	if err != nil {
		return domain.UserOverview{}, fmt.Errorf("get user overview by telegram_id: %w", err)
	}

	return overview, nil
}
