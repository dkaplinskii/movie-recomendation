package main

import (
	"context"
	"log"

	"studentgit.kata.academy/movie-recommendation-platform-20/movie-service/internal/app"
)

func main() {
	ctx := context.Background()

	if err := app.Start(ctx); err != nil {
		log.Fatalf("start app error: %v", err)
	}
}
