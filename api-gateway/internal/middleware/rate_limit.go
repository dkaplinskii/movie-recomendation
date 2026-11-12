package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type RateLimiter interface {
	Allow(ctx context.Context, ip string) (bool, error)
}

func RateLimitMiddleware(limiter RateLimiter, ignorePaths ...string) fiber.Handler {
	ignoreMap := make(map[string]bool, len(ignorePaths))
	for _, path := range ignorePaths {
		ignoreMap[path] = true
	}

	return func(c *fiber.Ctx) error {
		if ignoreMap[c.Path()] {
			return c.Next()
		}

		allowed, err := limiter.Allow(c.Context(), c.IP())
		if err != nil {
			// При ошибке Redis лучше пропустить запрос, чем блокировать всех
			return c.Next()
		}
		if !allowed {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "too many requests",
			})
		}
		return c.Next()
	}
}
