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
	fmt.Printf("Requesting URL: %s\n", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Failed to generate a request: %v\n", err)
		return "", err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		fmt.Printf("Failed to make a request: %v\n", err)
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return "", err
	}

	fmt.Printf("Response body: %s\n", string(body))

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Failed to get a successful response: %s\n", res.Status)
		return "", fmt.Errorf("non-200 response: %s", res.Status)
	}

	var krakenResponse KrakenResponse
	jsonErr := json.Unmarshal(body, &krakenResponse)
	if jsonErr != nil {
		fmt.Printf("Failed to unmarshal response: %v\n", jsonErr)
		return "", jsonErr
	}

	price, err := formatLastTradedPriceResponse(krakenResponse, pair)
	if err != nil {
		fmt.Printf("Failed to format response: %v\n", err)
		return "", err
	}

	// last traded price
	return price, nil
}

func formatLastTradedPriceResponse(krakenResponse KrakenResponse, pair string) (string, error) {
	ticker, ok := krakenResponse.Result[pair]
	if !ok || len(ticker.C) < 1 {
		fmt.Printf("Pair key not found or invalid format: %s\n", pair)
		return "", fmt.Errorf("invalid response format")
	}

	// format the ltp to two decimal places
	price, err := strconv.ParseFloat(ticker.C[0], 64)
	if err != nil {
		fmt.Printf("Failed to parse price: %v\n", err)
	}
	return fmt.Sprintf("%.2f", price), nil
}
