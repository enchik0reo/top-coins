package config

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/top-coins/internal/models"
)

type SenderType string

const (
	ApacheKafka SenderType = "kafka"
	LocalStdout SenderType = "local"
)

type FetcherConfig struct {
	Client HttpClient
	Kafka  KafkaProducer
	Sender SenderType `env:"FETCHER_SENDER" env-default:"local"`

	Env         models.Environment `env:"ENVIRONMENT" env-default:"local"`
	FetchTicker time.Duration      `env:"FETCH_TICKER" env-default:"60s"`
}

type KafkaProducer struct {
	KafkaBrokers []string      `env:"KAFKA_BROKERS" env-default:"localhost:9092"`
	KafkaTopic   string        `env:"KAFKA_TOPIC"`
	SourceName   string        `env:"FETCHER_NAME" env-default:"fetcher"`
	Timeot       time.Duration `env:"KAFKA_TIMEOUT" env-default:"10s"`
}

type HttpClient struct {
	Timeout time.Duration `env:"HTTP_TIMEOUT" env-default:"5s"`
}

func NewConfigFromEnv() (*FetcherConfig, error) {
	const fnName = "new config from env"

	cfg := FetcherConfig{}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, errors.Wrap(err, fnName)
	}

	return &cfg, nil
}
