package admins_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *AdminsRepository) AddAdmin(
	ctx context.Context,
	admin domain.Admin,
) (domain.Admin, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO vpnbot.admins (telegram_id, role, added_by, created_at)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (telegram_id)
	DO UPDATE 
	SET removed_at = NULL,
		role 	   = EXCLUDED.role,
		added_by   = EXCLUDED.added_by
	WHERE admins.removed_at IS NOT NULL
	RETURNING id, telegram_id, role, added_by, created_at;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		admin.TelegramID,
		string(admin.Role),
		admin.AddedBy,
		admin.CreatedAt,
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
				"admin with telegram_id=%d already active: %w",
				admin.TelegramID,
				core_errors.ErrConflict,
			)
		}
		return domain.Admin{}, fmt.Errorf("scan error: %w", err)
	}

	return adminModel.toDomain(), nil
}
