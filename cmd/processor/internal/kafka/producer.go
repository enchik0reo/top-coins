package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/segmentio/kafka-go"

	"github.com/top-coins/cmd/processor/internal/config"
	"github.com/top-coins/pkg/models"
)

type KafkaProducer struct {
	cfg    *config.Kafka
	writer *kafka.Writer
}

func newKafkaProducer(cfg *config.Kafka) *KafkaProducer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.CoinAlertTopic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{
		writer: w,
		cfg:    cfg,
	}
}

func (p *KafkaProducer) Publish(ctx context.Context, coinInfo models.CoinInfo) error {
	const fnName = "kafka producer: publish"

	ctx, cancel := context.WithTimeout(ctx, p.cfg.Timeot)
	defer cancel()

	keyBytes := []byte(coinInfo.Name)

	valueBytes, err := json.Marshal(coinInfo)
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
