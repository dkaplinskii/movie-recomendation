package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/config/modules"
)

func NewPool(ctx context.Context, postgresCfg modules.Postgres) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		postgresCfg.User,
		postgresCfg.Password,
		postgresCfg.Host,
		postgresCfg.Port,
		postgresCfg.DBName,
		postgresCfg.SSLMode,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return pool, nil
}
