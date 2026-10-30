package contract

import (
	"context"

	"github.com/top-coins/internal/models"
)

type CoinInfoFetcher interface {
	FetchInfo(ctx context.Context) (models.CoinInfo, error)
}
