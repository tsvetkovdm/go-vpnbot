package access_postgres_repository

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (r *AccessRepository) ListAccessibleConnections(
	ctx context.Context,
	userID int64,
) ([]domain.Connection, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT ` + connectionColumns + accessRule + `
	  AND u.id = $1
	ORDER BY c.sort_order ASC, c.id ASC;
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var connections []domain.Connection
	for rows.Next() {
		m, err := scanConnection(rows)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		connections = append(connections, m.toDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return connections, nil
}
