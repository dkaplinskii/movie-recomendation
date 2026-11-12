package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// HealthCheck godoc
// @Summary      Health check
// @Description  Check if the service is alive
// @Tags         System
// @Produce      plain
// @Success      200  {string}  string  "OK"
// @Router       /health [get]
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("OK")
}
