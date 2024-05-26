package handler

import (
	"github.com/antoniodaria/ltp_service/clients/kraken"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	KrakenClient kraken.Client
}

func NewHandler(app fiber.Router, krakenClient kraken.Client) *Handler {
	handler := &Handler{
		KrakenClient: krakenClient,
	}
	handler.connectRoutes(app)
	return handler
}

func (h *Handler) connectRoutes(app fiber.Router) {
	app.Get("/ltp", h.GetLTP)
}
