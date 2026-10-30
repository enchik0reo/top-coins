package sender

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/top-coins/cmd/fetcher/internal/config"
	"github.com/top-coins/cmd/fetcher/internal/sender/kafka"
	"github.com/top-coins/cmd/fetcher/internal/sender/local"
	"github.com/top-coins/internal/models"
)

type CoinInfoSender interface {
	Publish(ctx context.Context, coinsInfo []models.CoinInfo) error
	Close() error
}

type Sender struct {
	CoinInfoSender

	senderType config.SenderType
}

func NewInfoSender(conf *config.FetcherConfig) (CoinInfoSender, error) {
	const fnName = "new info sender"

	var cis CoinInfoSender

	switch conf.Sender {
	case config.ApacheKafka:
		sender, err := newKafkaSender()
		if err != nil {
			return nil, errors.Wrap(err, fnName)
		}

		cis = sender
	default:
		cis = newLocalSender()
	}

	return cis, nil
}

func newKafkaSender() (*kafka.KafkaProducer, error) {
	return nil, nil
}

func newLocalSender() *local.StdoutProducer {
	return local.NewStdoutProducer()
}
