package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/segmentio/kafka-go"

	"github.com/top-coins/cmd/fetcher/internal/config"
	"github.com/top-coins/pkg/logger"
	"github.com/top-coins/pkg/models"
)

var ErrEmptyBrokers = errors.New("empty brokers list")

type KafkaProducer struct {
	writer *kafka.Writer
	cfg    *config.KafkaProducer
	log    *logger.CustomLogger
}

func NewKafkaProducer(cfg *config.KafkaProducer, log *logger.CustomLogger) (*KafkaProducer, error) {
	const fnName = "new kafka producer"

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeot)
	defer cancel()

	if len(cfg.KafkaBrokers) == 0 {
		return nil, errors.Wrap(ErrEmptyBrokers, fnName)
	}

	brokerAddr := cfg.KafkaBrokers[0]
	conn, err := kafka.DialContext(ctx, "tcp", brokerAddr)
	if err != nil {
		return nil, errors.Wrap(err, fnName)
	}
	if err := conn.Close(); err != nil {
		return nil, errors.Wrap(err, fnName)
	}

	w := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.KafkaBrokers...),
		Topic:                  cfg.KafkaInfoTopic,
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

	p.log.Debug("Coin info published to kafka", p.log.Attr("func", fnName), p.log.Attr("count", len(coinsInfo)), p.log.Attr("topic", p.cfg.KafkaInfoTopic))

	return nil
}

func (p *KafkaProducer) Close() error {
	const fnName = "kafka producer: close"

	if err := p.writer.Close(); err != nil {
		return errors.Wrap(err, fnName)
	}

	return nil
}
