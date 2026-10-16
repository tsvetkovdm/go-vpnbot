package users_postgres_repository

import (
	"context"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (r *UsersRepository) GetUserOverviewByTelegramID(
	ctx context.Context,
	telegramID int64,
) (domain.UserOverview, error) {
	return r.getUserOverview(ctx, "u.telegram_id = $1", telegramID)
}
