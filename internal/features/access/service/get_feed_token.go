package access_service

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (s *AccessService) GetFeedToken(
	ctx context.Context,
	telegramID int64,
) (string, error) {
	token, err := s.accessRepository.EnsureFeedToken(ctx, telegramID, domain.NewFeedToken())
	if err != nil {
		return "", fmt.Errorf("get feed token: %w", err)
	}

	return token, nil
}
