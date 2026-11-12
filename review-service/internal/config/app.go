package config

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/config/modules"
)

type Config struct {
	AppName string `env:"APP_NAME" envDefault:"review-service"`
	AppEnv  string `env:"APP_ENV" envDefault:"local"`

	Log             modules.LoggerCfg `envPrefix:"APP_LOG_"`
	Postgres        modules.Postgres  `envPrefix:"POSTGRES_"`
	Kafka           modules.Kafka     `envPrefix:"KAFKA_"`
	Outbox          modules.Outbox    `envPrefix:"OUTBOX_"`
	GRPC            modules.GRPC      `envPrefix:"GRPC_"`
	ShutdownTimeout time.Duration     `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load env config: %w", err)
	}

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse env config: %w", err)
	}

	return cfg, nil
}
