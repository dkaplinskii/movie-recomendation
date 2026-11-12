package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"recommendation-service/internal/algorithm"
	"recommendation-service/internal/config"
	"recommendation-service/internal/controller/grpc"
	"recommendation-service/internal/infrastructure/stub"
	"recommendation-service/internal/logger"
	"recommendation-service/internal/observer"
	"recommendation-service/internal/service"
	"recommendation-service/internal/transport/grpcserver"
	"syscall"

	kafkainfra "recommendation-service/internal/infrastructure/kafka"

	redisinfra "recommendation-service/internal/infrastructure/redis"

	redisrepo "recommendation-service/internal/repository/redis"

	recommendationv1 "recommendation-service/pkg/recommendation/v1"
)

func AppRun(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log, err := logger.NewLogger(cfg.App.LogLevel, os.Stdout)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	ctx = logger.WithLogger(ctx, log)

	startupCtx, cancel := context.WithTimeout(ctx, cfg.Redis.RedisStartupTimeout)
	defer cancel()

	redisClient, err := redisinfra.NewRedis(startupCtx, cfg.Redis)
	if err != nil {
		return fmt.Errorf("init redis client: %w", err)
	}

	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Warnf("close redis client failed: %v", err)
		}
	}()

	recommendationRepository := redisrepo.NewRecommendationRepository(redisClient, cfg.RecommendationCache.RecommendationCacheTTL)

	syncProducer, err := kafkainfra.InitKafkaProducer(ctx, cfg.App.AppName, &cfg.Kafka)
	if err != nil {
		return fmt.Errorf("init kafka producer: %w", err)
	}
	kafkaProducer := kafkainfra.NewProducer(log, syncProducer)

	defer func() {
		if err := syncProducer.Close(); err != nil {
			log.Warnf("close kafka producer failed: %v", err)
		}
	}()

	consumerGroup, err := kafkainfra.InitKafkaConsumer(ctx, cfg.App.AppName, &cfg.Kafka)
	if err != nil {
		return fmt.Errorf("init kafka consumer: %w", err)
	}

	algorithmCfg := algorithm.Config(cfg.RecommendationAlgorithm)

	reviewClient := stub.NewReviewClient()
	movieClient := stub.NewMovieClient()
	recommendationService := service.NewService(recommendationRepository, log, kafkaProducer, cfg.Kafka.RecommendationEventsTopic, reviewClient, movieClient, algorithmCfg)

	reviewCreatedHandler := observer.NewReviewCreatedHandler(log, recommendationService)
	kafkaConsumer := observer.NewConsumer(log, reviewCreatedHandler)
	consumerRunner := kafkainfra.NewConsumerRunner(log, consumerGroup, cfg.Kafka.ReviewEventsTopic, kafkaConsumer)

	defer func() {
		if err := consumerRunner.Close(); err != nil {
			log.Warnf("close kafka consumer runner failed: %v", err)
		}
	}()

	server, err := grpcserver.NewServer(cfg.Server.GRPCPort, log)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	recommendationController := grpc.NewController(recommendationService, log)
	recommendationv1.RegisterRecommendationServiceServer(server.GRPCServer(), recommendationController)

	serverErrCh := make(chan error, 1)
	consumerErrCh := make(chan error, 1)

	go func() {
		if err := server.Start(); err != nil {
			serverErrCh <- err
		}
	}()

	go func() {
		if err := consumerRunner.Run(ctx); err != nil {
			consumerErrCh <- err
		}
	}()

	log.Info("App running")

	select {
	case <-ctx.Done():
		log.Info("shutdown signal received")
		server.Stop(cfg.Server.ShutdownTimeout)
		log.Info("application stopped")
		return nil

	case err := <-serverErrCh:
		server.Stop(cfg.Server.ShutdownTimeout)
		return fmt.Errorf("grpc server failed: %w", err)

	case err := <-consumerErrCh:
		log.Errorw("kafka consumer failed, stopping grpc server", "error", err)
		server.Stop(cfg.Server.ShutdownTimeout)
		return fmt.Errorf("kafka consumer failed: %w", err)
	}
}
