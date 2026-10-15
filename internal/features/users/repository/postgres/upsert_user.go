package users_postgres_repository

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (r *UsersRepository) UpsertUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO vpnbot.users (telegram_id, username, first_name, email, xray_uuid, status, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (telegram_id) DO UPDATE
	SET username   = EXCLUDED.username,
		first_name = EXCLUDED.first_name,
		updated_at  = EXCLUDED.updated_at,
		version    = users.version + 1
	RETURNING id, version, telegram_id, username, first_name, email, xray_uuid, status, created_at, updated_at;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.TelegramID,
		user.Username,
		user.FirstName,
		user.Email,
		user.XrayUUID,
		string(user.Status),
		user.CreatedAt,
		user.UpdatedAt,
	)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.TelegramID,
		&userModel.Username,
		&userModel.FirstName,
		&userModel.Email,
		&userModel.XrayUUID,
		&userModel.Status,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	return userModel.toDomain(), nil
}
