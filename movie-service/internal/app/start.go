package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"studentgit.kata.academy/movie-recommendation-platform-20/movie-service/internal/config"
	"studentgit.kata.academy/movie-recommendation-platform-20/movie-service/internal/infrastructure/postgres"
	"studentgit.kata.academy/movie-recommendation-platform-20/movie-service/internal/logger"
	"studentgit.kata.academy/movie-recommendation-platform-20/movie-service/internal/transport/grpcserver"
)

func Start(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	l, err := logger.New(cfg.Log.Level, os.Stdout)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	ctx = logger.WithLogger(ctx, l)
	defer func() {
		if err := l.Sync(); err != nil {
			log.Printf("logger sync error: %v\n", err)
		}
	}()

	if err := postgres.RunMigrations(ctx, cfg.Postgres); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	l.Info("migration completed")

	pgPool, err := postgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}
	defer pgPool.Close()
	l.Info("postgres connected")

	grpcServer, err := grpcserver.NewGRPCServer(cfg.GRPC.Port)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}
	l.Infow("grpc listener created", "port", cfg.GRPC.Port)

	if err := grpcServer.Run(ctx, cfg.ShutdownTimeout); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}

	l.Info("shutdown completed")

	return nil
}
