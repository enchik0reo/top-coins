package service

import (
	"context"
	"errors"
	"sync"
	"time"

	cockrErr "github.com/cockroachdb/errors"

	"github.com/top-coins/cmd/fetcher/internal/contract"
	"github.com/top-coins/internal/logger"
	"github.com/top-coins/internal/models"
)

type CoinInfoCollector struct {
	tickTimeout time.Duration
	fetchers    []contract.CoinInfoFetcher
	log         *logger.Logger
}

func NewCoinInfoCollector(ticker time.Duration, fetchers []contract.CoinInfoFetcher, logger *logger.Logger) *CoinInfoCollector {
	return &CoinInfoCollector{
		tickTimeout: ticker,
		fetchers:    fetchers,
		log:         logger,
	}
}

func (c *CoinInfoCollector) StartCollect(ctx context.Context) chan []models.CoinInfo {
	const fnName = "coin info collector: start"

	ticker := time.NewTicker(c.tickTimeout)

	infoChan := make(chan []models.CoinInfo)

	go func() {
		defer close(infoChan)
		defer ticker.Stop()

		resInfo, err := c.getCoinInfo(ctx)
		if err != nil {
			c.log.Error("failed to get coin info", "func", fnName, "error", err)

			select {
			case <-ctx.Done():
				return
			default:
			}
		}

		select {
		case <-ctx.Done():
			return
		case infoChan <- resInfo:
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				resInfo, err := c.getCoinInfo(ctx)
				if err != nil {
					c.log.Error("failed to get coin info", "func", fnName, "error", err)

					select {
					case <-ctx.Done():
						return
					default:
						continue
					}
				}

				select {
				case <-ctx.Done():
					return
				case infoChan <- resInfo:
				}
			}
		}
	}()

	return infoChan
}

func (c *CoinInfoCollector) getCoinInfo(ctx context.Context) ([]models.CoinInfo, error) {
	const fnName = "get coin info"

	ctx, cancel := context.WithTimeout(ctx, c.tickTimeout/2)
	defer cancel()

	resInfo := make([]models.CoinInfo, 0, len(c.fetchers))
	infoCh := make(chan models.CoinInfo)
	wg := sync.WaitGroup{}

	var err error
	mu := sync.Mutex{}

	for _, fetcher := range c.fetchers {
		wg.Go(func() {
			info, fetchErr := fetcher.FetchInfo(ctx)
			if fetchErr != nil {
				mu.Lock()
				err = errors.Join(err, fetchErr)
				mu.Unlock()

				return
			}

			infoCh <- info
		})
	}

	go func() {
		wg.Wait()
		close(infoCh)
	}()

	for info := range infoCh {
		resInfo = append(resInfo, info)
	}

	if len(resInfo) == 0 {
		return nil, cockrErr.Wrap(err, fnName)
	}

	return resInfo, nil
}
