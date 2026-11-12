package server

import (
	"time"

	"api-gateway/internal/config"
	"api-gateway/internal/infrastructure/log"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app             *fiber.App
	addr            string
	shutdownTimeout time.Duration
}

func New(cfg *config.Config, app *fiber.App) *Server {
	return &Server{
		app:             app,
		addr:            ":" + cfg.Port,
		shutdownTimeout: cfg.ShutdownTimeout,
	}
}

func (s *Server) Start() {
	go func() {
		log.Logger().Infof("Server starting on %s", s.addr)
		if err := s.app.Listen(s.addr); err != nil {
			log.Logger().Errorf("Server error: %v", err)
		}
	}()
}

func (s *Server) Stop() {
	log.Logger().Info("Shutting down server...")
	if err := s.app.ShutdownWithTimeout(s.shutdownTimeout); err != nil {
		log.Logger().Errorf("Server shutdown error: %v", err)
	}
	log.Logger().Info("Server exited")
}
