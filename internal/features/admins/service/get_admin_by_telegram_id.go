package admins_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *AdminsService) GetAdminByTelegramID(
	ctx context.Context,
	telegramID int64,
) (domain.Admin, error) {
	admin, err := s.adminsRepository.GetAdminByTelegramID(ctx, telegramID)
	if err != nil {
		return domain.Admin{}, fmt.Errorf("get admin from repository: %w", err)
	}

	return admin, nil
}
