package config

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Bot      BotConfig      `envPrefix:"BOT_"`
	Postgres PostgresConfig `envPrefix:"POSTGRES_"`
	Logger   LoggerConfig   `envPrefix:"LOGGER_"`
	Client   ClientConfig   `envPrefix:"CLIENT_"`
}

type BotConfig struct {
	Token   string        `env:"TOKEN,required"`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"10s"`
}

type PostgresConfig struct {
	Host     string `env:"HOST,required"`
	Port     int    `env:"PORT" envDefault:"5432"`
	User     string `env:"USER,required"`
	Password string `env:"PASSWORD,required"`
	Database string `env:"NAME,required"`
}

type LoggerConfig struct {
	Level string `env:"LEVEL" envDefault:"info"`
}

type ClientConfig struct {
	URL     string        `env:"URL,required"`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"10s"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error reading .env file: %w", err)
	}

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("config parsing error: %w", err)
	}

	return &cfg, nil
}

func (p *PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.User, p.Password, p.Host, p.Port, p.Database,
	)
}
