package access_transport_tg

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PublicURL    string `envconfig:"PUBLIC_URL" required:"true"`
	ProfileTitle string `envconfig:"PROFILE_TITLE" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("FEED", &config); err != nil {
		return Config{}, fmt.Errorf("envconfig process: %w", err)
	}

	if !strings.HasPrefix(config.PublicURL, "https://") {
		return Config{}, fmt.Errorf("FEED_PUBLIC_URL must start with https://, got %q", config.PublicURL)
	}
	config.PublicURL = strings.TrimSuffix(config.PublicURL, "/")

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get access TG config: %w", err)
		panic(err)
	}

	return config
}
