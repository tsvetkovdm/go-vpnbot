package payments_transport_tg

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	CardNumber string `envconfig:"CARD_NUMBER" required:"true"`
	SBPPhone   string `envconfig:"SBP_PHONE" required:"true"`
	Bank       string `envconfig:"BANK" required:"true"`
	Recipient  string `envconfig:"RECIPIENT" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("PAYMENT", &config); err != nil {
		return Config{}, fmt.Errorf("envconfig process: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get payments config: %w", err)
		panic(err)
	}

	return config
}
