package plans_postgres_repository

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

const planColumns = `id, version, title, is_active, created_at, updated_at`

const planPriceColumns = `id, version, plan_id, duration, price, is_active, created_at, updated_at`

type PlanModel struct {
	ID        int64
	Version   int64
	Title     string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m PlanModel) toDomain() domain.Plan {
	return domain.Plan{
		ID:        m.ID,
		Version:   m.Version,
		Title:     m.Title,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func scanPlan(row core_postgres_pool.Row) (PlanModel, error) {
	var m PlanModel
	err := row.Scan(
		&m.ID,
		&m.Version,
		&m.Title,
		&m.IsActive,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	return m, err
}

type PlanPriceModel struct {
	ID        int64
	Version   int64
	PlanID    int64
	Duration  pgtype.Interval
	Price     decimal.Decimal
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m PlanPriceModel) toDomain() domain.PlanPrice {
	return domain.PlanPrice{
		ID:        m.ID,
		Version:   m.Version,
		PlanID:    m.PlanID,
		Duration:  durationFromInterval(m.Duration),
		Price:     m.Price,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func scanPlanPrice(row core_postgres_pool.Row) (PlanPriceModel, error) {
	var m PlanPriceModel
	err := row.Scan(
		&m.ID,
		&m.Version,
		&m.PlanID,
		&m.Duration,
		&m.Price,
		&m.IsActive,
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
