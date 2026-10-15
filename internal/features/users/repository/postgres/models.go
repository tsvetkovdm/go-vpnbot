package users_postgres_repository

import (
	"time"

	"github.com/google/uuid"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

type UserModel struct {
	ID         int64
	Version    int64
	TelegramID int64
	Username   *string
	FirstName  string
	Email      string
	XrayUUID   uuid.UUID
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (m UserModel) toDomain() domain.User {
	return domain.User{
		ID:         m.ID,
		Version:    m.Version,
		TelegramID: m.TelegramID,
		Username:   m.Username,
		FirstName:  m.FirstName,
		Email:      m.Email,
		XrayUUID:   m.XrayUUID,
		Status:     domain.Status(m.Status),
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}
