package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port                    string        `env:"PORT" envDefault:"80"`
	ShutdownTimeout         time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
	DBHost                  string        `env:"DB_HOST" envDefault:"localhost"`
	DBPort                  string        `env:"DB_PORT" envDefault:"5432"`
	RedisAddr               string        `env:"REDIS_ADDR" envDefault:"redis:6379"`
	RateLimitExcludePaths   []string      `env:"RATE_LIMIT_EXCLUDE_PATHS" envSeparator:"," envDefault:"/health"`
	RateLimitPerMin         int           `env:"RATE_LIMIT_PER_MIN" envDefault:"100"`
	RateBurst               int           `env:"RATE_BURST" envDefault:"10"`
	ValidTokens             []string      `env:"VALID_TOKENS" envSeparator:"," envDefault:"secret-token-1,secret-token-2"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}
	return cfg, nil
}
