package access_transport_http

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ProfileTitle        string `envconfig:"PROFILE_TITLE" required:"true"`
	UpdateIntervalHours int    `envconfig:"UPDATE_INTERVAL_HOURS" required:"true"`
	SupportURL          string `envconfig:"SUPPORT_URL"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("FEED", &config); err != nil {
		return Config{}, fmt.Errorf("envconfig process: %w", err)
	}

	if config.UpdateIntervalHours < 1 {
		return Config{}, fmt.Errorf("FEED_UPDATE_INTERVAL_HOURS must be >= 1, got %d", config.UpdateIntervalHours)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get access HTTP config: %w", err)
		panic(err)
	}

	return config
}
