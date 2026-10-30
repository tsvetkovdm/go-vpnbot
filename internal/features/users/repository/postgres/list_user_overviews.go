package users_postgres_repository

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
)

func (r *UsersRepository) ListUserOverviews(
	ctx context.Context,
	filter domain.UserFilter,
	limit, offset int,
) ([]domain.UserOverview, int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	where, orderBy, err := filterSQL(filter)
	if err != nil {
		return nil, 0, fmt.Errorf("list user overviews: %w", err)
	}

	query := `
	SELECT ` + userOverviewColumns + `,
		   COUNT(*) OVER() AS total
	FROM vpnbot.users u
	LEFT JOIN vpnbot.subscriptions s ON s.user_id = u.id
	WHERE ` + where + `
	ORDER BY ` + orderBy + `
	LIMIT $1 OFFSET $2;
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var overviews []domain.UserOverview
	total := 0

	for rows.Next() {
		var m UserOverviewModel
		if err := rows.Scan(append(m.scanDest(), &total)...); err != nil {
			return nil, 0, fmt.Errorf("scan error: %w", err)
		}
		overviews = append(overviews, m.toDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("scan error: %w", err)
	}
	return overviews, total, nil
}

func filterSQL(filter domain.UserFilter) (where, orderBy string, err error) {
	switch filter {
	case domain.UserFilterAll:
		return "TRUE", "u.created_at DESC", nil
	case domain.UserFilterActive:
		return "u.status = 'active' AND s.expires_at > now()", "s.expires_at ASC", nil
	case domain.UserFilterExpiring:
		return "u.status = 'active' AND s.expires_at > now() AND s.expires_at <= now() + interval '7 days'", "s.expires_at ASC", nil
	case domain.UserFilterInactive:
		return "u.status = 'active' AND (s.id IS NULL OR s.expires_at <= now())", "u.created_at DESC", nil
	case domain.UserFilterBlocked:
		return "u.status = 'blocked'", "u.updated_at DESC", nil
	default:
		return "", "", fmt.Errorf("unknown filter %q: %w", filter, core_errors.ErrInvalidArgument)
	}
}
