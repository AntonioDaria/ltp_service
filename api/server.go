package api

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	app *fiber.App
}

func NewServer() *Server {
	return &Server{}
}

// StartAndListen
// @title Ascent API
// @version 1.0
// @description GITHub User Fav Language API
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
