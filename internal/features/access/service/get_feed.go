package access_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *AccessService) GetFeed(
	ctx context.Context,
	token string,
) (domain.Feed, error) {
	owner, err := s.accessRepository.GetFeedOwnerByToken(ctx, token)
	if err != nil {
		return domain.Feed{}, fmt.Errorf("get feed: %w", err)
	}

	connections, err := s.accessRepository.ListAccessibleConnections(ctx, owner.UserID)
	if err != nil {
		return domain.Feed{}, fmt.Errorf("get feed user_id=%d: %w", owner.UserID, err)
	}

	return domain.Feed{
		Owner:       owner,
		Connections: connections,
	}, nil
}
