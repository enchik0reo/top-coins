package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/segmentio/kafka-go"

	"github.com/top-coins/cmd/fetcher/internal/config"
	"github.com/top-coins/internal/models"
)

type KafkaProducer struct {
	writer *kafka.Writer
	cfg    *config.KafkaProducer
}

func NewKafkaProducer(cfg *config.KafkaProducer) *KafkaProducer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(cfg.KafkaBrokers...),
		Topic:    cfg.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{
		writer: w,
		cfg:    cfg,
	}
}

func (p *KafkaProducer) Publish(ctx context.Context, coinsInfo []models.CoinInfo) error {
	const fnName = "kafka producer: publish"

	ctx, cancel := context.WithTimeout(ctx, p.cfg.Timeot)
	defer cancel()

	keyBytes := []byte(coinsInfo[0].Name)

	valueBytes, err := json.Marshal(coinsInfo)
	if err != nil {
		return errors.Wrap(err, fnName)
	}

	message := kafka.Message{
		Key:   keyBytes,
		Value: valueBytes,
		Headers: []kafka.Header{
			{Key: "source", Value: []byte(p.cfg.SourceName)},
			{Key: "content-type", Value: []byte("application/json")},
		},
		Time: time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, message); err != nil {
		return errors.Wrap(err, fnName)
	}

	return nil
}

func (p *KafkaProducer) Close() error {
	const fnName = "kafka producer: close"

	if err := p.writer.Close(); err != nil {
		return errors.Wrap(err, fnName)
	}

	return nil
}
