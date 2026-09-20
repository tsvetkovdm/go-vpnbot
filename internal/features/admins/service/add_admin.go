package admins_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *AdminsService) AddAdmin(
	ctx context.Context,
	admin domain.Admin,
) (domain.Admin, error) {
	admin, err := s.adminsRepository.AddAdmin(ctx, admin)
	if err != nil {
		return domain.Admin{}, fmt.Errorf("add admin: %w", err)
	}

	return admin, err
}
