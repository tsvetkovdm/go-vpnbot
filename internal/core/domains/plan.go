package domain

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
)

const (
	PlanTitleMaxLen = 30
)

var (
	ErrPlanTitleEmpty   = fmt.Errorf("plan title is empty: %w", core_errors.ErrInvalidArgument)
	ErrPlanTitleTooLong = fmt.Errorf("plan title too long: %w", core_errors.ErrInvalidArgument)
)

type Plan struct {
	ID        int64
	Version   int64
	Title     string
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

func NewPlan(
	title string,
) (Plan, error) {
	title = strings.TrimSpace(title)

	if err := ValidatePlanTitle(title); err != nil {
		return Plan{}, err
	}

	return Plan{
		ID:       UninitializedID,
		Version:  UninitializedVersion,
		Title:    title,
		IsActive: false,
	}, nil
}
