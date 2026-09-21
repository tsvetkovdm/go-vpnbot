package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Plan struct {
	ID        int64
	Version   int64
	Title     string
	Price     decimal.Decimal
	Duration  Duration
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPlan(
	title string,
	price decimal.Decimal,
	duration Duration,
) Plan {
	now := time.Now()

	return Plan{
		ID:        UninitializedID,
		Version:   UninitializedVersion,
		Title:     title,
		Price:     price,
		Duration:  duration,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
