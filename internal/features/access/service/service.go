package access_service

import (
	"context"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

type AccessService struct {
	accessRepository AccessRepository
}

type AccessRepository interface {
	GetFeedOwnerByToken(
		ctx context.Context,
		token string,
	) (domain.FeedOwner, error)

	ListAccessibleConnections(
		ctx context.Context,
		userID int64,
	) ([]domain.Connection, error)
}

func NewAccessService(
	accessRepository AccessRepository,
) *AccessService {
	return &AccessService{
		accessRepository: accessRepository,
	}
}
