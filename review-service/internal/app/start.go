package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/config"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/controller"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/infrastructure/kafka"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/infrastructure/postgres"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/logger"
	outboxWorker "studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/outbox"
	kafkaPublisher "studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/publisher/kafka"
	outboxRepository "studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/repository/postgres/outbox"
	mainRepository "studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/repository/postgres/storage"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/service"
	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/internal/transport/grpcserver"
)

func Start(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(
		ctx,
		os.Interrupt,
		syscall.SIGTERM,
	)
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
			l.Errorw("logger sync error", "error", err)
		}
	}()

	if err := postgres.RunMigrations(ctx, cfg.Postgres); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	l.Info("migration completed")

	mainPool, err := postgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("connect main postgres pool: %w", err)
	}
	defer mainPool.Close()
	l.Info("main postgres pool connected")

	outboxPool, err := postgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("connect outbox postgres pool: %w", err)
	}
	defer outboxPool.Close()
	l.Info("outbox postgres pool connected")

	producer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		return fmt.Errorf("create kafka producer: %w", err)
	}
	defer producer.Close()
	l.Info("kafka connected")

	publisher := kafkaPublisher.NewPublisher(
		producer,
		cfg.Kafka.ReviewEventsTopic,
	)

	mainRepo := mainRepository.NewPostgresRepository(mainPool)

	outboxRepo := outboxRepository.NewPostgresOutboxRepository(
		outboxPool,
	)

	txManager := postgres.NewTransactionManager(mainPool)

	svc := service.NewService(
		mainRepo,
		outboxRepo,
		txManager,
	)

	ctrl := controller.NewController(svc)

	worker := outboxWorker.NewWorker(
		outboxRepo,
		publisher,
		cfg.Outbox,
	)

	grpcServer, err := grpcserver.NewGRPCServer(
		cfg.GRPC.Port,
		ctrl,
	)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}
	l.Infow("grpc listener created", "port", cfg.GRPC.Port)

	workerCtx, cancelWorker := context.WithCancel(ctx)
	workerDone := make(chan struct{})

	go func() {
		defer close(workerDone)
		worker.Run(workerCtx)
	}()

	runErr := grpcServer.Run(ctx, cfg.ShutdownTimeout)

	cancelWorker()
	<-workerDone

	if runErr != nil {
		return fmt.Errorf("run grpc server: %w", runErr)
	}

	l.Info("shutdown completed")

	return nil
}
