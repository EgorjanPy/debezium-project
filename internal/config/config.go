package config

import (
	"debez/pkg/postgres"
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string `env:"ENV" env-default:"development"`
	Server      Server
	Debezium    Debezium
	Postgres    postgres.Config
}
type Server struct {
	Port    int           `env:"PORT"         env-default:"8080"`
	TimeOut time.Duration `env:"HTTP_TIMEOUT" env-default:"30s"`
}
type Debezium struct {
	BaseURL string `env:"DEBEZIUM_BASE_URL" env-default:"http://localhost:8080"`
}

func ParseConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse confg: %w", err)
	}
	return cfg, nil
}
