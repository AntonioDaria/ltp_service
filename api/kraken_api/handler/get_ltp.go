package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type LastTradedPriceResponse struct {
	LTP []PairLTP `json:"ltp"`
}

type PairLTP struct {
	Pair   string `json:"pair"`
	Amount string `json:"amount"`
}

// GetLTP godoc
// @Summary Get Last Traded Price
// @Description Get the last traded price of a pair from Kraken api
// @Tags Kraken
// @Accept json
// @Produce json
// @Success 200 {object} LastTradedPriceResponse
// @Failure 500 {object} api.JSONError
// @Failure 405 {object} api.JSONError
// @Router /api/v1/ltp [get]
func (h *Handler) GetLTP(c *fiber.Ctx) error {
	pairs := map[string]string{
		"BTC/CHF": "XBTCHF",
		"BTC/EUR": "XXBTZEUR",
		"BTC/USD": "XXBTZUSD",
	}

	var response LastTradedPriceResponse

	for pair, krakenPair := range pairs {
		price, err := h.KrakenClient.GetLastTradedPrice(c.Context(), krakenPair)
		if err != nil {
			fmt.Printf("Failed to get last traded price for %s: %v", pair, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Error": "Failed to get last traded price",
			})
		}

		response.LTP = append(response.LTP, PairLTP{
			Pair:   pair,
			Amount: price,
		})
	}

	return c.JSON(response)
}
