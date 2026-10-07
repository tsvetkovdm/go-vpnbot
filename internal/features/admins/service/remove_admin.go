package admins_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
)

func (s *AdminsService) RemoveAdmin(
	ctx context.Context,
	telegramID int64,
) error {
	admin, err := s.GetAdminByTelegramID(ctx, telegramID)
	if err != nil {
		return fmt.Errorf("get admin: %w", err)
	}
	if admin.Role == domain.RoleSuper {
		return fmt.Errorf("admin telegram_id=%d is super: %w", telegramID, core_errors.ErrConflict)
	}

	if err := s.adminsRepository.RemoveAdmin(ctx, telegramID); err != nil {
		return fmt.Errorf("remove admin: %w", err)
	}

	return nil
}
