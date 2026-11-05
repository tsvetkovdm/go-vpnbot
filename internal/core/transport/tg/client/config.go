package core_tg_client

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Token     string `envconfig:"TOKEN" required:"true"`
	ProxyAddr string `envconfig:"PROXY_ADDR"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("TG", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig : %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get bot API config: %w", err)
		panic(err)
	}

	return config
}
