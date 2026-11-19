package kafka

import (
	"github.com/cockroachdb/errors"
	"github.com/top-coins/cmd/processor/internal/config"
	"github.com/top-coins/pkg/logger"
)

type KafkaService struct {
	cfg *config.Kafka
	log *logger.CustomLogger

	Consumer *KafkaConsumer
	Producer *KafkaProducer
}

func NewKafkaService(cfg *config.Kafka, log *logger.CustomLogger) (*KafkaService, error) {
	const fnName = "new kafka Service"

	producer, err := newKafkaProducer(cfg, log)
	if err != nil {
		return nil, errors.Wrap(err, fnName)
	}

	return &KafkaService{
		Consumer: newKafkaConsumer(cfg, log),
		Producer: producer,
		cfg:      cfg,
	}, nil
}
