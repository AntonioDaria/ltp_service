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
		if r.URL.String() != "/0/public/Ticker?pair=XXBTZUSD" {
			t.Errorf("Expected URL /0/public/Ticker?pair=XXBTZUSD, got %s", r.URL.String())
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
	price, err := client.GetLastTradedPrice(context.Background(), "XXBTZUSD")
	assert.NoError(t, err)

	// assert response
	assert.Equal(t, "1234.56", price)
}
