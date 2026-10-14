package domain

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/shopspring/decimal"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
)

const (
	PlanTitleMaxLen = 30
	PlanMaxMonths   = 24
	PlanMaxDays     = 366
)

var PlanMaxPrice = decimal.NewFromInt(1_000_000)

var (
	ErrPlanTitleEmpty      = fmt.Errorf("plan title is empty: %w", core_errors.ErrInvalidArgument)
	ErrPlanTitleTooLong    = fmt.Errorf("plan title too long: %w", core_errors.ErrInvalidArgument)
	ErrPlanDurationInvalid = fmt.Errorf("plan duration must be positive months or days, not both: %w", core_errors.ErrInvalidArgument)
	ErrPlanDurationTooLong = fmt.Errorf("plan duration too long: %w", core_errors.ErrInvalidArgument)
	ErrPlanPriceNegative   = fmt.Errorf("plan price is negative: %w", core_errors.ErrInvalidArgument)
	ErrPlanPriceTooHigh    = fmt.Errorf("plan price too high: %w", core_errors.ErrInvalidArgument)
	ErrPlanPricePrecision  = fmt.Errorf("plan price precision is invalid: %w", core_errors.ErrInvalidArgument)
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

func ValidatePlanTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrPlanTitleEmpty
	}

	if n := utf8.RuneCountInString(title); n > PlanTitleMaxLen {
		return fmt.Errorf("%w: got %d, max %d", ErrPlanTitleTooLong, n, PlanTitleMaxLen)
	}

	return nil
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

func NewPlan(
	title string,
	duration Duration,
	price decimal.Decimal,
) (Plan, error) {
	title = strings.TrimSpace(title)

	if err := ValidatePlanTitle(title); err != nil {
		return Plan{}, err
	}
	if err := ValidatePlanDuration(duration); err != nil {
		return Plan{}, err
	}
	if err := ValidatePlanPrice(price); err != nil {
		return Plan{}, err
	}

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
	}, nil
}
