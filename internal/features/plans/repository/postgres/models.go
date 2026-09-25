package plans_postgres_repository

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

type PlanModel struct {
	ID        int64
	Version   int64
	Title     string
	Duration  pgtype.Interval
	Price     decimal.Decimal
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m PlanModel) toDomain() domain.Plan {
	return domain.Plan{
		ID:      m.ID,
		Version: m.Version,
		Title:   m.Title,
		Duration: domain.NewDuration(
			m.Duration.Months,
			m.Duration.Days,
			m.Duration.Microseconds,
		),
		Price:     m.Price,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
