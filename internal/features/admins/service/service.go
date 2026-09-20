package admins_service

import (
	"context"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

type AdminsService struct {
	adminsRepository AdminsRepository
}

type AdminsRepository interface {
	GetAdminByTelegramID(
		ctx context.Context,
		telegramID int64,
	) (domain.Admin, error)

	EnsureSuperAdmin(
		ctx context.Context,
		admin domain.Admin,
	) error

	AddAdmin(
		ctx context.Context,
		admin domain.Admin,
	) (domain.Admin, error)

	RemoveAdmin(
		ctx context.Context,
		telegramID int64,
	) error
}

func NewAdminsService(
	adminsRepository AdminsRepository,
) *AdminsService {
	return &AdminsService{
		adminsRepository: adminsRepository,
	}
}
