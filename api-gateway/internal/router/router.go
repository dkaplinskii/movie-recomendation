package router

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"api-gateway/internal/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
)

// Router структура для хранения маршрутов
type Router struct {
	app *fiber.App
}

// New создаёт новый экземпляр Router
func New() *Router {
	return &Router{
		app: fiber.New(),
	}
}

// SetupRoutes настраивает все маршруты
func (r *Router) SetupRoutes(cfg *config.Config) {
	// Публичные маршруты
	r.app.Get("/health", handlers.HealthCheck)

	// Swagger UI
	r.app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Защищённые маршруты (с аутентификацией)
	api := r.app.Group("/api", middleware.AuthMiddleware(cfg.ValidTokens))
	// Здесь добавляй защищённые маршруты, например:
	// api.Get("/users", handlers.GetUsers)
	// api.Post("/movies", handlers.CreateMovie)
	_ = api // Временно игнорируем, чтобы не было ошибки "declared and not used"
}

// GetApp возвращает экземпляр fiber.App
func (r *Router) GetApp() *fiber.App {
	return r.app
}
