package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
)

func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	telegramID int64,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM vpnbot.users
	WHERE telegram_id=$1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, telegramID)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"user with telegram_id='%d': %w",
			telegramID,
			core_errors.ErrNotFound,
		)
	}

	return nil
}
