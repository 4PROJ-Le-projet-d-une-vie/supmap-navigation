package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	APIServerHost string `env:"API_SERVER_HOST"`
	APIServerPort string `env:"API_SERVER_PORT"`
	RedisHost     string `env:"REDIS_HOST"`
	RedisPort     string `env:"REDIS_PORT"`
}

func New() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &cfg, nil
}
