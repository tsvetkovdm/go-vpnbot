package domain

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
)

const (
	PlanMaxMonths = 24
	PlanMaxDays   = 366
)

var PlanMaxPrice = decimal.NewFromInt(1_000_000)

var (
	ErrPlanDurationInvalid = fmt.Errorf("plan duration must be positive months or days, not both: %w", core_errors.ErrInvalidArgument)
	ErrPlanDurationTooLong = fmt.Errorf("plan duration too long: %w", core_errors.ErrInvalidArgument)
	ErrPlanPriceNegative   = fmt.Errorf("plan price is negative: %w", core_errors.ErrInvalidArgument)
	ErrPlanPriceTooHigh    = fmt.Errorf("plan price too high: %w", core_errors.ErrInvalidArgument)
	ErrPlanPricePrecision  = fmt.Errorf("plan price precision is invalid: %w", core_errors.ErrInvalidArgument)
)

type PlanPrice struct {
	ID        int64
	Version   int64
	PlanID    int64
	Duration  Duration
	Price     decimal.Decimal
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPlanPrice(
	planID int64,
	duration Duration,
	price decimal.Decimal,
) (PlanPrice, error) {
	if err := ValidatePlanDuration(duration); err != nil {
		return PlanPrice{}, err
	}
	if err := ValidatePlanPrice(price); err != nil {
		return PlanPrice{}, err
	}

	return PlanPrice{
		ID:       UninitializedID,
		Version:  UninitializedVersion,
		PlanID:   planID,
		Duration: duration,
		Price:    price,
		IsActive: true,
	}, nil
}

func ValidatePlanDuration(d Duration) error {
	if d.Months < 0 || d.Days < 0 || d.Microseconds != 0 {
		return fmt.Errorf("%w: %+v", ErrPlanDurationInvalid, d)
	}

	if (d.Months > 0) == (d.Days > 0) {
		return fmt.Errorf("%w: %+v", ErrPlanDurationInvalid, d)
	}

	if d.Months > PlanMaxMonths || d.Days > PlanMaxDays {
		return fmt.Errorf("%w: %+v", ErrPlanDurationTooLong, d)
	}

	return nil
}

func ValidatePlanPrice(price decimal.Decimal) error {
	if price.IsNegative() {
		return fmt.Errorf("%w: %s", ErrPlanPriceNegative, price)
	}

	if price.GreaterThan(PlanMaxPrice) {
		return fmt.Errorf("%w: %s", ErrPlanPriceTooHigh, price)
	}

	if !price.Equal(price.Round(2)) {
		return fmt.Errorf("%w: %s", ErrPlanPricePrecision, price)
	}

	return nil
}
