package admins_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *AdminsService) EnsureSuperAdmin(
	ctx context.Context,
	telegramID int64,
) error {
	admin := domain.NewAdmin(telegramID, domain.RoleSuper, nil)

	if err := s.adminsRepository.EnsureSuperAdmin(ctx, admin); err != nil {
		return fmt.Errorf("ensure super admin: %w", err)
	}

	return nil
}
