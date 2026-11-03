package access_service

import (
	"context"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_xray "github.com/tsvetkovdm/go-vpnbot/internal/core/xray"
)

type AccessService struct {
	accessRepository AccessRepository
	xrayClient       XrayClient
}

type AccessRepository interface {
	GetFeedOwnerByToken(
		ctx context.Context,
		token string,
	) (domain.FeedOwner, error)

	EnsureFeedToken(
		ctx context.Context,
		telegramID int64,
		newToken string,
	) (string, error)

	ListAccessibleConnections(
		ctx context.Context,
		userID int64,
	) ([]domain.Connection, error)

	ListInboundTags(
		ctx context.Context,
	) ([]string, error)

	ListDesiredAccess(
		ctx context.Context,
	) ([]domain.AccessEntry, error)
}

type XrayClient interface {
	ListUsers(ctx context.Context, inboundTag string) ([]core_xray.User, error)
	AddUser(ctx context.Context, inboundTag string, user core_xray.User) error
	RemoveUser(ctx context.Context, inboundTag string, email string) error
}

func NewAccessService(
	accessRepository AccessRepository,
	xrayClient XrayClient,
) *AccessService {
	return &AccessService{
		accessRepository: accessRepository,
		xrayClient:       xrayClient,
	}
}
