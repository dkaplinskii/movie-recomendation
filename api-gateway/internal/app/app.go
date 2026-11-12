package app

import (
	"context"
	"fmt"

	"api-gateway/internal/config"
	"api-gateway/internal/ratelimiter"
	"api-gateway/internal/router"
	"api-gateway/internal/server"
)

func Run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config load: %w", err)
	}

	rateLimiter := ratelimiter.NewIPRateLimiter(cfg.RateLimitPerMin, cfg.RateBurst, cfg.RedisAddr)
	_ = rateLimiter

	r := router.New()
	r.SetupRoutes(cfg)

	srv := server.New(cfg, r.GetApp())
	srv.Start()

	return nil
}
