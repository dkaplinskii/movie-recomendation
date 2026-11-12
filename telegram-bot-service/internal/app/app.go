package app

import (
	"context"
	"fmt"

	"telegram-bot-service/internal/clients/gateway"
	"telegram-bot-service/internal/config"
	"telegram-bot-service/internal/logger"
	"telegram-bot-service/internal/migrator"
	"telegram-bot-service/internal/service"
	"telegram-bot-service/internal/transport/telegram"

	infraPostgres "telegram-bot-service/internal/infrastructure/postgres"
	storagePostgres "telegram-bot-service/internal/storage/postgres"
)

func Run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	log, err := logger.InitLogger(cfg.Logger.Level)
	if err != nil {
		return fmt.Errorf("error logger init: %w", err)
	}
	defer log.Sync()

	log.Info("Start telegram-bot service...")

	if err := migrator.NewMigrationPostgres(cfg.Postgres.DSN()); err != nil {
		return fmt.Errorf("database migrations failed: %w", err)
	}
	log.Info("Migrations applied successfully")

	pool, err := infraPostgres.NewPool(ctx, cfg.Postgres.DSN())
	if err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}
	defer pool.Close()

	userRepo := storagePostgres.NewUserRepository(pool)
	client := gateway.NewClient(cfg.Client.URL, cfg.Client.Timeout)
	botService := service.NewBotService(userRepo, client)

	tgBot, err := telegram.NewBot(cfg.Bot.Token, cfg.Bot.Timeout, log, botService)
	if err != nil {
		return fmt.Errorf("error creating bot: %w", err)
	}

	go func() {
		<-ctx.Done()
		log.Info("Stop signal received...")
		tgBot.Stop()
	}()

	log.Info("Bot launched successfully")
	tgBot.Start()

	log.Info("Service stopped")
	return nil
}
