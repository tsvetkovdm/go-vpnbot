package domain

import (
	"fmt"
	"time"

	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
)

const (
	SubscriptionMaxExtendDays = 366
	SubscriptionMaxYearsAhead = 5
)

var (
	ErrSubscriptionDaysInvalid = fmt.Errorf("extend days out of range: %w", core_errors.ErrInvalidArgument)
	ErrSubscriptionDateInPast  = fmt.Errorf("paid until date in past: %w", core_errors.ErrInvalidArgument)
	ErrSubscriptionDateTooFar  = fmt.Errorf("paid until date too far: %w", core_errors.ErrInvalidArgument)
)

type Subscription struct {
	ID             int64
	Version        int64
	UserID         int64
	ExpiresAt      time.Time
	AnchorDay      int
	Notified3DayAt *time.Time
	Notified1DayAt *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewSubscription(
	userID int64,
	now time.Time,
) Subscription {
	today := dateOf(now)

	return Subscription{
		ID:        UninitializedID,
		Version:   UninitializedVersion,
		UserID:    userID,
		ExpiresAt: today,
		AnchorDay: today.Day(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (s Subscription) HasAccess(now time.Time) bool {
	return now.Before(s.ExpiresAt)
}

func (s Subscription) PaidUntil(loc *time.Location) time.Time {
	return dateOf(s.ExpiresAt.In(loc)).AddDate(0, 0, -1)
}

func (s Subscription) Extend(d Duration, now time.Time) Subscription {
	today := dateOf(now)
	paid := s.PaidUntil(now.Location())
	anchor := s.AnchorDay

	if paid.Before(today) {
		paid = today
		anchor = today.Day()
	}

	if d.Months > 0 {
		paid = addMonths(paid, int(d.Months), anchor)
	}

	if d.Days > 0 {
		paid = paid.AddDate(0, 0, int(d.Days))
		anchor = paid.Day()
	}

	s.ExpiresAt = paid.AddDate(0, 0, 1)
	s.AnchorDay = anchor
	s.Notified3DayAt, s.Notified1DayAt = nil, nil
	s.UpdatedAt = now

	return s
}

func (s Subscription) SetPaidUntil(date, now time.Time) (Subscription, error) {
	today := dateOf(now)
	date = dateOf(date.In(now.Location()))
	if date.Before(today) {
		return Subscription{}, ErrSubscriptionDateInPast
	}

	if date.After(today.AddDate(SubscriptionMaxYearsAhead, 0, 0)) {
		return Subscription{}, ErrSubscriptionDateTooFar
	}

	s.ExpiresAt = date.AddDate(0, 0, 1)
	s.AnchorDay = date.Day()
	s.Notified3DayAt, s.Notified1DayAt = nil, nil
	s.UpdatedAt = now

	return s, nil
}

func ValidateExtendDays(days int) error {
	if days < 1 || days > SubscriptionMaxExtendDays {
		return ErrSubscriptionDaysInvalid
	}

	return nil
}

func dateOf(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func daysInMonth(year int, month time.Month, loc *time.Location) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, loc).Day()
}

func addMonths(date time.Time, n, anchor int) time.Time {
	first := time.Date(date.Year(), date.Month()+time.Month(1), 1, 0, 0, 0, 0, date.Location())
	day := min(anchor, daysInMonth(first.Year(), first.Month(), date.Location()))

	return time.Date(first.Year(), first.Month(), day, 0, 0, 0, 0, date.Location())
}
