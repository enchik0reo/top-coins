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

const (
	coinbaseURL = "https://api.coinbase.com/v2/prices/BTC-USD/spot" // {"data":{"amount":"79744.025","base":"BTC","currency":"USD"}}
	price       = "amount"
)

type coinbaseCoinInfo struct {
	Data map[string]string `json:"data"`
}

type CoinbaseClient struct {
	*CoinsClient
	url string
	log *logger.CustomLogger
}

func newCoinbaseClient(cl *CoinsClient) *CoinbaseClient {
	bc := CoinbaseClient{
		CoinsClient: cl,
		url:         coinbaseURL,
		log:         cl.log,
	}

	return &bc
}

func (c *CoinbaseClient) FetchInfo(ctx context.Context) (models.CoinInfo, error) {
	const fnName = "coinbase client: fetch info"

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
			c.log.Error("Fail to close response body", c.log.Attr("error", err), c.log.Attr("fnName", fnName))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return models.CoinInfo{}, errors.Wrap(ErrorStatusCode, fnName)
	}

	info := coinbaseCoinInfo{}

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	coinInfo, err := c.mapCoin(info)
	if err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	return coinInfo, nil
}

func (c *CoinbaseClient) mapCoin(info coinbaseCoinInfo) (models.CoinInfo, error) {
	const fnName = "map coin"

	btcPriceStr, ok := info.Data[price]
	if !ok {
		return models.CoinInfo{}, errors.Wrap(ErrorEmptyPrice, fnName)
	}

	f, err := strconv.ParseFloat(btcPriceStr, 64)
	if err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	cents := uint64(f * 100)

	return models.CoinInfo{
		Name:      models.BitCoin,
		Price:     cents,
		Currency:  models.CentsUSD,
		Source:    models.Coinbase,
		CreatedAt: time.Now().UTC(),
	}, nil
}
