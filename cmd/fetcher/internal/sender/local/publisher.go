package local

import (
	"context"
	"fmt"
	"os"

	"github.com/top-coins/pkg/logger"
	"github.com/top-coins/pkg/models"
)

type StdoutProducer struct {
	log *logger.CustomLogger
}

func NewStdoutProducer(log *logger.CustomLogger) *StdoutProducer {
	return &StdoutProducer{log: log}
}

func (p *StdoutProducer) Publish(ctx context.Context, coinsInfo []models.CoinInfo) error {
	const fnName = "stdout producer: publish"

	fmt.Fprintf(os.Stdout, "Coins:\n%+v\n", coinsInfo)
	p.log.Debug("Coin info published to stdout", p.log.Attr("func", fnName), p.log.Attr("count", len(coinsInfo)))

	return nil
}

func (p *StdoutProducer) Close() error {
	return nil
}
