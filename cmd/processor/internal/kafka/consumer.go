package kafka

import (
	"context"
	"encoding/json"

	"github.com/cockroachdb/errors"
	"github.com/segmentio/kafka-go"
	"github.com/top-coins/cmd/processor/internal/config"
	"github.com/top-coins/pkg/logger"
	"github.com/top-coins/pkg/models"
)

var ErrEmptyInfo = errors.New("empty information field")

type KafkaConsumer struct {
	cfg    *config.Kafka
	log    *logger.CustomLogger
	reader *kafka.Reader
}

func newKafkaConsumer(cfg *config.Kafka, log *logger.CustomLogger) *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          cfg.CoinInfoTopic,
		GroupID:        cfg.ConsumerGroupID,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		MaxWait:        cfg.Timeot,
		CommitInterval: 0,
		StartOffset:    kafka.LastOffset,
	})

	return &KafkaConsumer{
		cfg:    cfg,
		log:    log,
		reader: r,
	}
}

func (c *KafkaConsumer) StartConsume(ctx context.Context) (<-chan *CoinMessage, error) {
	const fnName = "kafka consumer: consume"
	outChan := make(chan *CoinMessage, 100)

	go func() {
		defer close(outChan)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := c.reader.FetchMessage(ctx)
				if err != nil {
					if errors.Is(err, context.Canceled) {
						return
					}
					c.log.Error("Can't fetch kafka message", c.log.Attr("error", err), c.log.Attr("func", fnName))

					continue
				}

				coinMsg := c.mustProcessMessage(msg)

				select {
				case <-ctx.Done():
					return
				case outChan <- coinMsg:
				}
			}
		}
	}()

	return outChan, nil
}

func (c *KafkaConsumer) mustProcessMessage(msg kafka.Message) *CoinMessage {
	const fnName = "must process message "

	var coinsInfo []models.CoinInfo
	if err := json.Unmarshal(msg.Value, &coinsInfo); err != nil {
		c.log.Error("Can't unmarshal kafka message value", c.log.Attr("error", err), c.log.Attr("func", fnName))

		return &CoinMessage{
			info: coinsInfo,
			err:  err,
			commitFunc: func(cCtx context.Context) error {
				return c.reader.CommitMessages(cCtx, msg)
			},
		}
	}

	if len(coinsInfo) == 0 {
		return &CoinMessage{
			info: coinsInfo,
			err:  ErrEmptyInfo,
			commitFunc: func(cCtx context.Context) error {
				return c.reader.CommitMessages(cCtx, msg)
			},
		}
	}

	return &CoinMessage{
		info: coinsInfo,
		commitFunc: func(cCtx context.Context) error {
			return c.reader.CommitMessages(cCtx, msg)
		},
	}
}

func (c *KafkaConsumer) Close() error {
	const fnName = "kafka consumer: close"

	if err := c.reader.Close(); err != nil {
		return errors.Wrap(err, fnName)
	}

	return nil
}
