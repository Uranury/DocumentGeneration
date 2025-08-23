package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ListenAddr string `env:"LISTEN_ADDR" default:":8080"`
	JWTSecret  string `env:"JWT_SECRET" required:"true"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &cfg, nil
}
