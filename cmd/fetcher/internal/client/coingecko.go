package client

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/top-coins/pkg/logger"
	"github.com/top-coins/pkg/models"
)

const (
	coingeckoURL = "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd" // {"bitcoin":{"usd":79782}}
	usd          = "usd"
)

type coingeckoCoinInfo struct {
	BTCData map[string]uint64 `json:"bitcoin"`
}

type CoingeckoClient struct {
	*CoinsClient
	url string
	log *logger.Logger
}

func newCoingeckoClient(cl *CoinsClient) *CoingeckoClient {
	bc := CoingeckoClient{
		CoinsClient: cl,
		url:         coingeckoURL,
		log:         cl.log,
	}

	return &bc
}

func (c *CoingeckoClient) FetchInfo(ctx context.Context) (models.CoinInfo, error) {
	const fnName = "Coingecko client: fetch info"

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

	info := coingeckoCoinInfo{}

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	coinInfo, err := c.mapCoin(info)
	if err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	return coinInfo, nil
}

func (c *CoingeckoClient) mapCoin(info coingeckoCoinInfo) (models.CoinInfo, error) {
	const fnName = "map coin"

	btcPrice, ok := info.BTCData[usd]
	if !ok {
		return models.CoinInfo{}, errors.Wrap(ErrorEmptyPrice, fnName)
	}

	cents := uint64(btcPrice * 100)

	return models.CoinInfo{
		Name:      models.BitCoin,
		Price:     cents,
		Currency:  models.CentsUSD,
		Source:    models.Coingecko,
		CreatedAt: time.Now().UTC(),
	}, nil
}
