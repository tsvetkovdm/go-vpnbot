package admins_postgres_repository

import (
	"context"
	"fmt"
	"time"

	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
)

func (r *AdminsRepository) RemoveAdmin(
	ctx context.Context,
	telegramID int64,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE vpnbot.admins
	SET removed_at = $2
	WHERE telegram_id = $1 AND removed_at IS NULL;
	`

	cmdTag, err := r.pool.Exec(ctx, query, telegramID, time.Now())
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"admin with telegram_id=%d: %w",
			telegramID,
			core_errors.ErrNotFound,
		)
	}

	return nil
}
