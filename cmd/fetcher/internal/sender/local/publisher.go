package local

import (
	"context"
	"fmt"
	"os"

	"github.com/top-coins/pkg/models"
)

type StdoutProducer struct {
}

func NewStdoutProducer() *StdoutProducer {
	return &StdoutProducer{}
}

func (p *StdoutProducer) Publish(ctx context.Context, coinsInfo []models.CoinInfo) error {
	fmt.Fprintf(os.Stdout, "Coins:\n%+v\n", coinsInfo)
	return nil
}

func (p *StdoutProducer) Close() error {
	return nil
}
