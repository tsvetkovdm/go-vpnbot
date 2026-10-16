package users_postgres_repository

import (
	"context"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (r *UsersRepository) GetUserOverviewByUsername(
	ctx context.Context,
	username string,
) (domain.UserOverview, error) {
	return r.getUserOverview(ctx, "lower(u.username) = lower($1)", username)
}
