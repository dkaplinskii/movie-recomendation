package middleware

import (
	"log"
	"strings"
	"github.com/gofiber/fiber/v2"
)

var publicPaths = []string{
	"/health",
	"/swagger",
	"/metrics",
	"/favicon.ico",
}

func AuthMiddleware(validTokens []string) fiber.Handler {
	tokenMap := make(map[string]bool, len(validTokens))
	for _, token := range validTokens {
		tokenMap[token] = true
	}

	return func(c *fiber.Ctx) error {
		path := c.Path()

		if isPublicPath(path) {
			return c.Next()
		}

		token := c.Get("x-auth-token")
		if token == "" || !tokenMap[token] {
			log.Printf("[WARNING] Unauthorized attempt: %s %s from %s",
				c.Method(), path, c.IP())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "unauthorized",
				"message": "Invalid or missing authentication token",
			})
		}

		return c.Next()
	}
}

func isPublicPath(path string) bool {
	for _, public := range publicPaths {
		if path == public || strings.HasPrefix(path, public+"/") {
			return true
		}
	}
	return false
}
