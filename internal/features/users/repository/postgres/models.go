package users_postgres_repository

import (
	"time"

	"github.com/google/uuid"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

const userOverviewColumns = `
	u.id, u.version, u.telegram_id, u.username, u.first_name, u.email,
	u.xray_uuid, u.status, u.created_at, u.updated_at,
	s.id, s.version, s.expires_at, s.anchor_day,
	EXISTS (
		SELECT 1 FROM vpnbot.admins a
		WHERE a.telegram_id = u.telegram_id AND a.removed_at IS NULL
	),
	(
		SELECT pay.status FROM vpnbot.payments pay
		WHERE pay.user_id = u.id AND pay.status IN ('created', 'pending')
	),
	(
		SELECT last.plan_title_snapshot FROM vpnbot.payments last
		WHERE last.user_id = u.id AND last.status = 'successful'
		ORDER BY last.resolved_at DESC, last.id DESC
		LIMIT 1
	)`

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

type UserOverviewModel struct {
	User         UserModel
	SubID        *int64
	SubVersion   *int64
	SubExpiresAt *time.Time
	SubAnchorDay *int16
	IsAdmin      bool

	UnresolvedPaymentStatus *string
	LastPlanTitle           *string
}

func (m UserOverviewModel) toDomain() domain.UserOverview {
	overview := domain.UserOverview{
		User:          m.User.toDomain(),
		IsAdmin:       m.IsAdmin,
		LastPlanTitle: m.LastPlanTitle,
	}

	if m.UnresolvedPaymentStatus != nil {
		status := domain.PaymentStatus(*m.UnresolvedPaymentStatus)
		overview.UnresolvedPaymentStatus = &status
	}

	if m.SubID != nil {
		overview.Subscription = &domain.Subscription{
			ID:        *m.SubID,
			Version:   *m.SubVersion,
			UserID:    m.User.ID,
			ExpiresAt: *m.SubExpiresAt,
			AnchorDay: int(*m.SubAnchorDay),
		}
	}

	return overview
}

func (m *UserOverviewModel) scanDest() []any {
	return []any{
		&m.User.ID,
		&m.User.Version,
		&m.User.TelegramID,
		&m.User.Username,
		&m.User.FirstName,
		&m.User.Email,
		&m.User.XrayUUID,
		&m.User.Status,
		&m.User.CreatedAt,
		&m.User.UpdatedAt,
		&m.SubID,
		&m.SubVersion,
		&m.SubExpiresAt,
		&m.SubAnchorDay,
		&m.IsAdmin,
		&m.UnresolvedPaymentStatus,
		&m.LastPlanTitle,
	}
}
