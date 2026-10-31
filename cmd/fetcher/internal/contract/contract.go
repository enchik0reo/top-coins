package contract

import (
	"context"

	"github.com/top-coins/pkg/models"
)

type CoinInfoFetcher interface {
	FetchInfo(ctx context.Context) (models.CoinInfo, error)
}
