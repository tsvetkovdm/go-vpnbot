package admins_postgres_repository

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
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
	RETURNING id, telegram_id, role, added_by, created_at
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
		return domain.Admin{}, fmt.Errorf("scan error: %w", err)
	}

	return adminModel.toDomain(), nil
}
