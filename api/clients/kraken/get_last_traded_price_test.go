package kraken

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLastTradedPrice(t *testing.T) {
	// create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// write response
		response := KrakenResponse{
			Result: map[string]KrakenTicker{
				"XXBTZUSD": {
					C: []string{"1234.56"},
				},
			},
		}

		// check request method
		if r.Method != http.MethodGet {
			t.Errorf("Expected method %s, got %s", http.MethodGet, r.Method)
		}

		// check request URL
		if r.URL.Path != "/0/public/Ticker" {
			t.Errorf("Expected URL %s, got %s", "/0/public/Ticker", r.URL.Path)
		}

		// send response to be tested
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	// create client
	client := NewClient(ts.Client(), ts.URL)

	// test GetLastTradedPrice
	t.Run("Valid Response", func(t *testing.T) {
		// test GetLastTradedPrice
		price, err := client.GetLastTradedPrice(context.Background(), "XXBTZUSD")
		assert.NoError(t, err)

		// assert response
		assert.Equal(t, "1234.56", price)
	})

	t.Run("Invalid Pair", func(t *testing.T) {
		// test GetLastTradedPrice with an invalid pair
		_, err := client.GetLastTradedPrice(context.Background(), "INVALID")
		assert.Error(t, err)
	})

	t.Run("Empty Response", func(t *testing.T) {
		// test server with empty response
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := KrakenResponse{
				Result: map[string]KrakenTicker{},
			}
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				t.Fatal(err)
			}
		}))
		defer ts.Close()

		client := NewClient(ts.Client(), ts.URL)

		// test GetLastTradedPrice with empty response
		_, err := client.GetLastTradedPrice(context.Background(), "XXBTZUSD")
		assert.Error(t, err)
	})

	t.Run("Network Error", func(t *testing.T) {
		// create client with a non-existent server
		client := NewClient(&http.Client{}, "http://nonexistent")

		// test GetLastTradedPrice with a network error
		_, err := client.GetLastTradedPrice(context.Background(), "XXBTZUSD")
		assert.Error(t, err)
	})
}
