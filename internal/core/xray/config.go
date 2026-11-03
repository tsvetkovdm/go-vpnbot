package core_xray

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	APIAddr string        `envconfig:"API_ADDR" required:"true"`
	Timeout time.Duration `envconfig:"TIMEOUT" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("XRAY", &config); err != nil {
		return Config{}, fmt.Errorf("envconfig process: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get xray config: %w", err)
		panic(err)
	}

	return config
}
