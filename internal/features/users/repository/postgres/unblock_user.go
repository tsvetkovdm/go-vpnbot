package users_postgres_repository

import (
	"context"
	"fmt"
	"time"

	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
)

func (r *UsersRepository) UnblockUser(
	ctx context.Context,
	telegramID int64,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE vpnbot.users
	SET status = 'active', updated_at = $2, version = version + 1
	WHERE telegram_id=$1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, telegramID, time.Now())
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"user with telegram_id=%d: %w",
			telegramID,
			core_errors.ErrNotFound,
		)
	}

	return nil
}
