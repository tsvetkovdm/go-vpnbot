package core_config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
)

type Config struct {
	TimeZone             *time.Location
	SuperAdminTelegramID int64
	AdminContact         string
}

func NewConfig() (*Config, error) {
	tz := os.Getenv("TIME_ZONE")
	if tz == "" {
		tz = "UTC"
	}

	zone, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("load time zone: %s: %w", tz, err)
	}

	arg := os.Getenv("TG_SUPER_ADMIN_ID")
	telegramID, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse telegram_id %q: %w", arg, core_errors.ErrInvalidArgument)
	}

	contact := os.Getenv("TG_ADMIN_CONTACT")
	if contact == "" {
		return nil, fmt.Errorf("TG_ADMIN_CONTACT is empty: %w", core_errors.ErrInvalidArgument)
	}

	return &Config{
		TimeZone:             zone,
		SuperAdminTelegramID: telegramID,
		AdminContact:         contact,
	}, nil
}

func NewConfigMust() *Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get core config: %w", err)
		panic(err)
	}

	return config
}
