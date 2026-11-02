package plans_postgres_repository

import (
	"context"
	"fmt"
)

func (r *PlansRepository) CanUseFreePrices(
	ctx context.Context,
	telegramID int64,
) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT NOT EXISTS (
		SELECT 1
		FROM vpnbot.subscriptions s
		JOIN vpnbot.users u ON u.id = s.user_id
		WHERE u.telegram_id = $1
	);
	`

	var canUse bool
	if err := r.pool.QueryRow(ctx, query, telegramID).Scan(&canUse); err != nil {
		return false, fmt.Errorf("scan error: %w", err)
	}

	return canUse, nil
}
