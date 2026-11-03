package access_postgres_repository

import (
	"context"
	"fmt"
)

func (r *AccessRepository) ListInboundTags(
	ctx context.Context,
) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT DISTINCT c.inbound_tag
	FROM vpnbot.connections c
	ORDER BY c.inbound_tag;
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return tags, nil
}
