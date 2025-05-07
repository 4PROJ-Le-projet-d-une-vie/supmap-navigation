package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Env string

const (
	EnvProd Env = "prod"
	EnvDev  Env = "dev"
)

func (e Env) IsValid() bool {
	switch e {
	case EnvProd, EnvDev:
		return true
	}
	return false
}

type Config struct {
	APIServerHost         string `env:"API_SERVER_HOST"`
	APIServerPort         string `env:"API_SERVER_PORT"`
	RedisHost             string `env:"REDIS_HOST"`
	RedisPort             string `env:"REDIS_PORT"`
	RedisIncidentsChannel string `env:"REDIS_INCIDENTS_CHANNEL"`
	Env                   Env    `env:"ENV" envDefault:"prod"`
}

func New() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if !cfg.Env.IsValid() {
		return nil, fmt.Errorf("invalid env variable (must be 'prod' or 'dev')")
	}
	return &cfg, nil
}
