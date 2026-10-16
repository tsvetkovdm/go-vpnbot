package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) GetUserOverview(
	ctx context.Context,
	userID int64,
) (domain.UserOverview, error) {
	return r.getUserOverview(ctx, "u.id = $1", userID)
}

func (r *UsersRepository) getUserOverview(
	ctx context.Context,
	where string,
	arg any,
) (domain.UserOverview, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT ` + userOverviewColumns + `
	FROM vpnbot.users u
	LEFT JOIN vpnbot.subscriptions s ON s.user_id = u.id
	WHERE ` + where + `
	ORDER BY u.updated_at DESC
	LIMIT 1;
	`

	var m UserOverviewModel
	if err := r.pool.QueryRow(ctx, query, arg).Scan(m.scanDest()...); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.UserOverview{}, fmt.Errorf(
				"user overview where %s (%v): %w", where, arg, core_errors.ErrNotFound,
			)
		}
		return domain.UserOverview{}, fmt.Errorf("scan error: %w", err)
	}

	return m.toDomain(), nil
}