package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"telegram-bot-service/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil {
		log.Fatalf("app start error: %v", err)
	}
}
