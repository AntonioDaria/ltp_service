package kraken

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type KrakenResponse struct {
	Result map[string]KrakenTicker `json:"result"`
}

type KrakenTicker struct {
	C []string `json:"c"` // closing price
}

func (c *ClientImplementation) GetLastTradedPrice(ctx context.Context, pair string) (string, error) {
	url := fmt.Sprintf("%s/0/public/Ticker?pair=%s", c.BaseUrl, pair)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Failed to generate a request %v", err)
		return "", err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		fmt.Printf("Failed to make a request %v", err)
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Failed to read response body %v", err)
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Failed to get a successful response %v", err)
		return "", err
	}

	var krakenResponse KrakenResponse
	jsonErr := json.Unmarshal(body, &krakenResponse)
	if jsonErr != nil {
		fmt.Printf("Failed to unmarshal response %v", jsonErr)
		return "", jsonErr
	}

	price, err := formatLastTradedPriceResponse(krakenResponse, pair)
	if err != nil {
		return "", err
	}

	// last traded price
	return price, nil
}

func formatLastTradedPriceResponse(krakenResponse KrakenResponse, pair string) (string, error) {
	ticker, ok := krakenResponse.Result[pair]
	if !ok || len(ticker.C) < 1 {
		return "", fmt.Errorf("invalid response format")
	}

	// format the ltp to two decimal places
	price, err := strconv.ParseFloat(ticker.C[0], 64)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.2f", price), nil
}
