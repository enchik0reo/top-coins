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
	krakenURL   = "https://api.kraken.com/0/public/Ticker?pair=XBTUSD" // {"error":[],"result":{"XXBTZUSD":{"a":["79737.70000","3","3.000"],"b":["79737.60000","6","6.000"],"c":["79737.70000","0.01686558"],"v":["632.80885430","2375.08376615"],"p":["79750.75772","79938.16212"],"t":[52251,131712],"l":["79378.50000","79335.70000"],"h":["80238.20000","80701.70000"],"o":"80001.80000"}}}
	tradingPair = "XXBTZUSD"
)

type krakenCoinInfo struct {
	Result map[string]ticker `json:"result"`
	Errors []error           `json:"error"`
}

type ticker struct {
	Ask          []string `json:"a"`
	Bid          []string `json:"b"`
	Close        []string `json:"c"`
	Volume       []string `json:"v"`
	VWAP         []string `json:"p"`
	Trades       []int    `json:"t"`
	Low          []string `json:"l"`
	High         []string `json:"h"`
	OpeningPrice string   `json:"o"`
}

type KrakenClient struct {
	*CoinsClient
	url string
	log *logger.CustomLogger
}

func newKrakenClient(cl *CoinsClient) *KrakenClient {
	bc := KrakenClient{
		CoinsClient: cl,
		url:         krakenURL,
		log:         cl.log,
	}

	return &bc
}

func (c *KrakenClient) FetchInfo(ctx context.Context) (models.CoinInfo, error) {
	const fnName = "kraken client: fetch info"

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

	info := krakenCoinInfo{}

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	if len(info.Errors) > 0 {
		return models.CoinInfo{}, errors.Wrap(errors.Join(info.Errors...), fnName)
	}

	coinInfo, err := c.mapCoin(info)
	if err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	return coinInfo, nil
}

func (c *KrakenClient) mapCoin(info krakenCoinInfo) (models.CoinInfo, error) {
	const fnName = "map coin"

	btcData, ok := info.Result[tradingPair]
	if !ok || len(btcData.Close) == 0 {
		return models.CoinInfo{}, errors.Wrap(ErrorEmptyPrice, fnName)
	}

	f, err := strconv.ParseFloat(btcData.Close[0], 64)
	if err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	cents := uint64(f * 100)

	return models.CoinInfo{
		Name:      models.BitCoin,
		Price:     cents,
		Currency:  models.CentsUSD,
		Source:    models.Kraken,
		CreatedAt: time.Now().UTC(),
	}, nil
}
