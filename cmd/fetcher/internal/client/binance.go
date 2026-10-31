package client

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/top-coins/pkg/logger"
	"github.com/top-coins/pkg/models"
)

const binanceURL = "https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT" // {"symbol":"BTCUSDT","price":"79900.01000000"}

type binanceCoinInfo struct {
	CoinName string `json:"symbol"`
	PriceStr string `json:"price"`
}

type BinanceClient struct {
	*CoinsClient
	url string
	log *logger.Logger
}

func newBinanceClient(cl *CoinsClient) *BinanceClient {
	bc := BinanceClient{
		CoinsClient: cl,
		url:         binanceURL,
		log:         cl.log,
	}

	return &bc
}

func (c *BinanceClient) FetchInfo(ctx context.Context) (models.CoinInfo, error) {
	const fnName = "binance client: fetch info"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.log.Error("fail to close response body", "error", err, "fnName", fnName)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return models.CoinInfo{}, errors.Wrap(ErrorStatusCode, fnName)
	}

	info := binanceCoinInfo{}

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	coinInfo, err := c.mapCoin(info)
	if err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	return coinInfo, nil
}

func (c *BinanceClient) mapCoin(info binanceCoinInfo) (models.CoinInfo, error) {
	const fnName = "map coin"

	f, err := strconv.ParseFloat(info.PriceStr, 64)
	if err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	cents := uint64(f * 100)

	return models.CoinInfo{
		Name:      models.BitCoin,
		Price:     cents,
		Currency:  models.CentsUSD,
		Source:    models.Binance,
		CreatedAt: time.Now().UTC(),
	}, nil
}
