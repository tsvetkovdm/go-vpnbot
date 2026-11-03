package access_transport_worker

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Interval time.Duration `envconfig:"INTERVAL" required:"true"`
	DryRun   bool          `envconfig:"DRY_RUN" default:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("RECONCILE", &config); err != nil {
		return Config{}, fmt.Errorf("envconfig process: %w", err)
	}

	if config.Interval < 10*time.Second {
		return Config{}, fmt.Errorf("RECONCILE_INTERVAL must be >= 10s, got %s", config.Interval)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get reconcile worker config: %w", err)
		panic(err)
	}

	return config
}
