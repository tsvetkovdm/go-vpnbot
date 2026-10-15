package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) GetUserByTelegramID(
	ctx context.Context,
	telegramID int64,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, telegram_id, username, first_name, email, xray_uuid, status, created_at, updated_at
	FROM vpnbot.users
	WHERE telegram_id=$1;
	`

	row := r.pool.QueryRow(ctx, query, telegramID)

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
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with telegram_id: %d: %w",
				telegramID,
				core_errors.ErrNotFound,
			)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	return userModel.toDomain(), nil
}
