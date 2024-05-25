package api

import (
	"context"
	"fmt"

	"github.com/antoniodaria/ltp_service/api/clients/kraken"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	app          *fiber.App
	krakenClient kraken.Client
}

func NewServer(krakenClient kraken.Client) *Server {
	return &Server{
		krakenClient: krakenClient,
	}
}

// StartAndListen
// @title Last Traded Price API
// @version 1.0
// @description API for retrieval of the Last Traded Price of Bitcoin
// @host localhost:8000
// @BasePath /
func (s *Server) StartAndListen(ctx context.Context) {
	// Fiber
	app := fiber.New()
	s.app = app

	// This will recover from panics anywhere in the stack
	app.Use(recover.New(
		recover.Config{
			EnableStackTrace: true,
		},
	))
	app.Use(cors.New())

	// Add health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	err := app.Listen(fmt.Sprintf("0.0.0.0:%d", 8000)) // TODO: add port to config s.Config.APIPort
	if err != nil {
		fmt.Printf("Failed to start/stop API: %v", err)
	}
}

func (s *Server) Shutdown() {
	if s.app != nil {
		_ = s.app.Shutdown()
	}
}
