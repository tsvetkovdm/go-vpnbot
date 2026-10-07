package admins_postgres_repository

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (r *AdminsRepository) EnsureSuperAdmin(
	ctx context.Context,
	admin domain.Admin,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO vpnbot.admins (telegram_id, role, added_by, created_at)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (telegram_id) 
	DO UPDATE 
	SET role = EXCLUDED.role, removed_at = NULL;
	`

	if _, err := r.pool.Exec(
		ctx,
		query,
		admin.TelegramID,
		string(admin.Role),
		admin.AddedBy,
		admin.CreatedAt,
	); err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	return nil
}
