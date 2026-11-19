package sender

import (
	"context"

	"github.com/cockroachdb/errors"

	"github.com/top-coins/cmd/fetcher/internal/config"
	"github.com/top-coins/cmd/fetcher/internal/sender/kafka"
	"github.com/top-coins/cmd/fetcher/internal/sender/local"
	"github.com/top-coins/pkg/logger"
	"github.com/top-coins/pkg/models"
)

type CoinInfoSender interface {
	Publish(ctx context.Context, coinsInfo []models.CoinInfo) error
	Close() error
}

type Sender struct {
	CoinInfoSender

	senderType config.SenderType
}

func NewInfoSender(conf *config.FetcherConfig, log *logger.CustomLogger) (CoinInfoSender, error) {
	const fnName = "new info sender"

	var cis CoinInfoSender

	switch conf.Sender {
	case config.ApacheKafka:
		sender, err := newKafkaSender(&conf.Kafka, log)
		if err != nil {
			return nil, errors.Wrap(err, fnName)
		}

		cis = sender
	default:
		cis = newLocalSender(log)
	}

	log.Debug("Info sender created", log.Attr("func", fnName), log.Attr("sender", conf.Sender))

	return cis, nil
}

func newKafkaSender(cfg *config.KafkaProducer, log *logger.CustomLogger) (*kafka.KafkaProducer, error) {
	const fnName = "new kafka sender"

	producer, err := kafka.NewKafkaProducer(cfg, log)
	if err != nil {
		return nil, errors.Wrap(err, fnName)
	}

	return producer, nil
}

func newLocalSender(log *logger.CustomLogger) *local.StdoutProducer {
	return local.NewStdoutProducer(log)
}
