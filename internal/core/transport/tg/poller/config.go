package core_tg_poller

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PollTimeout     time.Duration `envconfig:"POLL_TIMEOUT" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("TG", &config); err != nil {
		return Config{}, fmt.Errorf("envconfig process: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get TG poller config: %w", err)
		panic(err)
	}

	return config
}
