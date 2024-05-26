package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_Kraken "github.com/antoniodaria/ltp_service/clients/kraken/mock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_GetLTPClientFailure(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockKrakenClient := mock_Kraken.NewMockClient(ctrl)

	mockKrakenClient.EXPECT().GetLastTradedPrice(gomock.Any(), gomock.Any()).Return("", assert.AnError)

	handler := &Handler{KrakenClient: mockKrakenClient}

	app := fiber.New()
	app.Get("/api/v1/ltp", handler.GetLTP)

	// Test
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ltp", nil)
	resp, _ := app.Test(req)
	defer resp.Body.Close()

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_GetLTPSuccess(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockKrakenClient := mock_Kraken.NewMockClient(ctrl)

	app := fiber.New()

	handler := &Handler{KrakenClient: mockKrakenClient}
	app.Get("/api/v1/ltp", handler.GetLTP)

	// Expectations
	mockKrakenClient.EXPECT().GetLastTradedPrice(gomock.Any(), "XBTCHF").Return("49000.12", nil)
	mockKrakenClient.EXPECT().GetLastTradedPrice(gomock.Any(), "XXBTZEUR").Return("50000.12", nil)
	mockKrakenClient.EXPECT().GetLastTradedPrice(gomock.Any(), "XXBTZUSD").Return("52000.12", nil)

	// Test
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ltp", nil)
	resp, _ := app.Test(req)
	defer resp.Body.Close()

	// Convert response body to byte
	resBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	expectedResponse := `{"ltp":[{"pair":"BTC/CHF","amount":"49000.12"},{"pair":"BTC/EUR","amount":"50000.12"},{"pair":"BTC/USD","amount":"52000.12"}]}`
	assert.JSONEq(t, expectedResponse, string(resBody))
}

func Test_GetLTPPartialFailure(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockKrakenClient := mock_Kraken.NewMockClient(ctrl)

	app := fiber.New()

	handler := &Handler{KrakenClient: mockKrakenClient}
	app.Get("/api/v1/ltp", handler.GetLTP)

	// Expectations
	mockKrakenClient.EXPECT().GetLastTradedPrice(gomock.Any(), "XBTCHF").Return("49000.12", nil)
	mockKrakenClient.EXPECT().GetLastTradedPrice(gomock.Any(), "XXBTZEUR").Return("", fmt.Errorf("API error"))

	// Test
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ltp", nil)
	resp, _ := app.Test(req)
	defer resp.Body.Close()

	// Convert response body to byte
	resBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	expectedResponse := `{"Error": "Failed to get last traded price"}`
	assert.JSONEq(t, expectedResponse, string(resBody))
}

func Test_GetLTPInvalidMethod(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockKrakenClient := mock_Kraken.NewMockClient(ctrl)

	app := fiber.New()

	handler := &Handler{KrakenClient: mockKrakenClient}
	app.Get("/api/v1/ltp", handler.GetLTP)

	// Test with POST method
	req := httptest.NewRequest(http.MethodPost, "/api/v1/ltp", nil)
	resp, _ := app.Test(req)
	defer resp.Body.Close()

	// Assertions
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func Test_GetLTPNoPairsAvailable(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockKrakenClient := mock_Kraken.NewMockClient(ctrl)

	app := fiber.New()

	// Mock handler with no pairs
	_ = &Handler{KrakenClient: mockKrakenClient}
	app.Get("/api/v1/ltp", func(c *fiber.Ctx) error {
		return c.JSON(LastTradedPriceResponse{LTP: []PairLTP{}})
	})

	// Test
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ltp", nil)
	resp, _ := app.Test(req)
	defer resp.Body.Close()

	// Convert response body to byte
	resBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	expectedResponse := `{"ltp":[]}`
	assert.JSONEq(t, expectedResponse, string(resBody))
}
