package admins_postgres_repository

import (
	"context"
	"fmt"
)

func (r *AdminsRepository) CountPendingPayments(
	ctx context.Context,
) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT COUNT(*)
	FROM vpnbot.payments
	WHERE status = 'pending';
	`

	var count int
	if err := r.pool.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("scan error: %w", err)
	}

	return count, nil
}
