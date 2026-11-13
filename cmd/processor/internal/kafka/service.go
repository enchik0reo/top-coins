package kafka

import (
	"github.com/top-coins/cmd/processor/internal/config"
	"github.com/top-coins/pkg/logger"
)

type KafkaService struct {
	cfg *config.Kafka
	log *logger.CustomLogger

	Consumer *KafkaConsumer
	Producer *KafkaProducer
}

func NewKafkaService(cfg *config.Kafka, log *logger.CustomLogger) *KafkaService {
	return &KafkaService{
		Consumer: newKafkaConsumer(cfg, log),
		Producer: newKafkaProducer(cfg),
		cfg:      cfg,
	}
}
