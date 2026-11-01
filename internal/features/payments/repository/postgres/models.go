package payments_postgres_repository

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

const paymentColumns = `
	p.id, p.user_id, p.plan_price_id,
	p.plan_title_snapshot, p.price_snapshot, p.duration_snapshot,
	p.comment_code, p.status, p.resolved_by_admin_id, p.resolved_at,
	p.created_at, p.updated_at`

type PaymentModel struct {
	ID                int64
	UserID            int64
	PlanPriceID       int64
	PlanTitle         string
	Price             decimal.Decimal
	Duration          pgtype.Interval
	CommentCode       int16
	Status            string
	ResolvedByAdminID *int64
	ResolvedAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (m PaymentModel) toDomain() domain.Payment {
	return domain.Payment{
		ID:                m.ID,
		UserID:            m.UserID,
		PlanPriceID:       m.PlanPriceID,
		PlanTitle:         m.PlanTitle,
		Price:             m.Price,
		Duration:          durationFromInterval(m.Duration),
		CommentCode:       int(m.CommentCode),
		Status:            domain.PaymentStatus(m.Status),
		ResolvedByAdminID: m.ResolvedByAdminID,
		ResolvedAt:        m.ResolvedAt,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func scanPayment(row core_postgres_pool.Row) (PaymentModel, error) {
	var m PaymentModel
	err := row.Scan(
		&m.ID,
		&m.UserID,
		&m.PlanPriceID,
		&m.PlanTitle,
		&m.Price,
		&m.Duration,
		&m.CommentCode,
		&m.Status,
		&m.ResolvedByAdminID,
		&m.ResolvedAt,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	return m, err
}

func durationFromInterval(i pgtype.Interval) domain.Duration {
	return domain.NewDuration(i.Months, i.Days, i.Microseconds)
}

func intervalFromDuration(d domain.Duration) pgtype.Interval {
	return pgtype.Interval{
		Months:       d.Months,
		Days:         d.Days,
		Microseconds: d.Microseconds,
		Valid:        true, // false означало бы NULL
	}
}
