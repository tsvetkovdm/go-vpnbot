package admins_postgres_repository

import (
	"time"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

type AdminModel struct {
	ID         int64
	TelegramID int64
	Role       string
	AddedBy    *int64
	CreatedAt  time.Time
}

func (m AdminModel) toDomain() domain.Admin {
	return domain.Admin{
		ID:         m.ID,
		TelegramID: m.TelegramID,
		Role:       domain.Role(m.Role),
		AddedBy:    m.AddedBy,
		CreatedAt:  m.CreatedAt,
	}
}
