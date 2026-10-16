package subscriptions_postgres_repository

import (
	"time"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

const subscriptionColumns = `id, version, user_id, expires_at, anchor_day,
	notified_3_day_at, notified_1_day_at, created_at, updated_at`

type SubscriptionModel struct {
	ID             int64
	Version        int64
	UserID         int64
	ExpiresAt      time.Time
	AnchorDay      int16
	Notified3DayAt *time.Time
	Notified1DayAt *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (m SubscriptionModel) toDomain() domain.Subscription {
	return domain.Subscription{
		ID:             m.ID,
		Version:        m.Version,
		UserID:         m.UserID,
		ExpiresAt:      m.ExpiresAt,
		AnchorDay:      int(m.AnchorDay),
		Notified3DayAt: m.Notified3DayAt,
		Notified1DayAt: m.Notified1DayAt,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func scanSubscription(row core_postgres_pool.Row) (SubscriptionModel, error) {
	var m SubscriptionModel
	err := row.Scan(
		&m.ID,
		&m.Version,
		&m.UserID,
		&m.ExpiresAt,
		&m.AnchorDay,
		&m.Notified3DayAt,
		&m.Notified1DayAt,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	return m, err
}
