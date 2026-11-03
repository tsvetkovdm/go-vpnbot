package access_postgres_repository

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (r *AccessRepository) ListDesiredAccess(
	ctx context.Context,
) ([]domain.AccessEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT DISTINCT c.inbound_tag, u.email, u.xray_uuid` + accessRule + `
	ORDER BY c.inbound_tag, u.email;
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var entries []domain.AccessEntry
	for rows.Next() {
		var e domain.AccessEntry
		if err := rows.Scan(&e.InboundTag, &e.Email, &e.XrayUUID); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		entries = append(entries, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return entries, nil
}
