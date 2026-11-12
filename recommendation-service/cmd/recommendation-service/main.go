package main

import (
	"context"
	"log"

	"recommendation-service/internal/app"
)

func main() {
	ctx := context.Background()

	if err := app.AppRun(ctx); err != nil {
		log.Fatalf("start app error: %v", err)
	}
}
