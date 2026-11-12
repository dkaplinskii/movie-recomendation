package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"studentgit.kata.academy/movie-recommendation-platform-20/movie-service/internal/config/modules"
)

func RunMigrations(ctx context.Context, postgresCfg modules.Postgres) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		postgresCfg.User,
		postgresCfg.Password,
		postgresCfg.Host,
		postgresCfg.Port,
		postgresCfg.DBName,
		postgresCfg.SSLMode,
	)

	postgresDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("migration error: %w", err)
	}
	defer postgresDB.Close()

	if err := postgresDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping postgres for migrations: %w", err)
	}
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("error with set dialect: %w", err)
	}
	if err := goose.Up(postgresDB, postgresCfg.MigrationsPath); err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	return nil
}
