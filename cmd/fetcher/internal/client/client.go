package client

import (
	"errors"
	"net/http"
	"time"

	"github.com/top-coins/cmd/fetcher/internal/contract"
	"github.com/top-coins/pkg/logger"
)

var (
	ErrorStatusCode = errors.New("unexpected response status code")
	ErrorEmptyPrice = errors.New("empty price")
)

type CoinsClient struct {
	timeout    time.Duration
	httpClient http.Client
	log        *logger.Logger
}

func NewCoinsClient(tm time.Duration, log *logger.Logger) *CoinsClient {
	cc := CoinsClient{
		timeout: tm,
	}

	cc.httpClient = http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     5 * tm,
		},
		Timeout: tm,
	}

	return &cc
}

func (c *CoinsClient) BringFetchers() []contract.CoinInfoFetcher {
	return []contract.CoinInfoFetcher{newBinanceClient(c),
		newCoinbaseClient(c),
		newCoingeckoClient(c),
		newKrakenClient(c),
	}
}
