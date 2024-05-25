package kraken

import (
	"context"
	"net/http"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE

type Client interface {
	GetLastTradedPrice(ctx context.Context, pair string) (string, error)
}

type ClientImplementation struct {
	HTTPClient *http.Client
	BaseUrl    string
}

func NewClient(httpClient *http.Client, baseUrl string) *ClientImplementation {
	return &ClientImplementation{
		HTTPClient: httpClient,
		BaseUrl:    baseUrl,
	}
}
