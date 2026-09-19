package admins_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *AdminsRepository) GetAdminByTelegramID(
	ctx context.Context,
	telegramID int64,
) (domain.Admin, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, telegram_id, role, added_by, created_at
	FROM vpnbot.admins
	WHERE telegram_id=$1;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		telegramID,
	)

	var adminModel AdminModel
	err := row.Scan(
		&adminModel.ID,
		&adminModel.TelegramID,
		&adminModel.Role,
		&adminModel.AddedBy,
		&adminModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Admin{}, fmt.Errorf(
				"admin with telegram_id=%d: %w",
				telegramID,
				core_errors.ErrNotFound,
			)
		}
		return domain.Admin{}, fmt.Errorf("scan error: %w", err)
	}

	return adminModel.toDomain(), nil
}
