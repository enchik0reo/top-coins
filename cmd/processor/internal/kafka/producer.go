package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/segmentio/kafka-go"

	"github.com/top-coins/cmd/processor/internal/config"
	"github.com/top-coins/pkg/logger"
	"github.com/top-coins/pkg/models"
)

var ErrEmptyBrokers = errors.New("empty brokers list")

type KafkaProducer struct {
	cfg    *config.Kafka
	log    *logger.CustomLogger
	writer *kafka.Writer
}

func newKafkaProducer(cfg *config.Kafka, log *logger.CustomLogger) (*KafkaProducer, error) {
	const fnName = "new kafka producer"

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeot)
	defer cancel()

	if len(cfg.Brokers) == 0 {
		return nil, errors.Wrap(ErrEmptyBrokers, fnName)
	}

	brokerAddr := cfg.Brokers[0]
	conn, err := kafka.DialContext(ctx, "tcp", brokerAddr)
	if err != nil {
		return nil, errors.Wrap(err, fnName)
	}
	if err := conn.Close(); err != nil {
		return nil, errors.Wrap(err, fnName)
	}

	w := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Brokers...),
		Topic:                  cfg.CoinAlertTopic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
		MaxAttempts:            20,
	}

	return &KafkaProducer{
		writer: w,
		cfg:    cfg,
		log:    log,
	}, nil
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

	p.log.Debug("Alert published to kafka", p.log.Attr("func", fnName), p.log.Attr("coin", coinInfo.Name), p.log.Attr("topic", p.cfg.CoinAlertTopic))

	return nil
}

func (p *KafkaProducer) Close() error {
	const fnName = "kafka producer: close"

	if err := p.writer.Close(); err != nil {
		return errors.Wrap(err, fnName)
	}

	return nil
}
