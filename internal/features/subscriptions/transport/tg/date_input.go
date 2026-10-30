package subscriptions_transport_tg

import (
	"errors"
	"strings"
	"time"
)

var errDateFormat = errors.New("invalid date format")

func parseDateInput(input string, now time.Time) (time.Time, error) {
	s := strings.TrimSpace(input)
	loc := now.Location()

	if d, err := time.ParseInLocation("2.1.2006", s, loc); err == nil {
		return d, nil
	}

	if d, err := time.ParseInLocation("2.1", s, loc); err == nil {
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		date := time.Date(now.Year(), d.Month(), d.Day(), 0, 0, 0, 0, loc)
		if date.Before(today) {
			date = date.AddDate(1, 0, 0)
		}
		return date, nil
	}

	return time.Time{}, errDateFormat
}
