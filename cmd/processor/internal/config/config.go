package config

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/top-coins/pkg/models"
)

type SenderType string

const (
	ApacheKafka SenderType = "kafka"
	LocalStdout SenderType = "local"
)

type Kafka struct {
	Brokers         []string      `env:"KAFKA_BROKERS" env-default:"localhost:9092"`
	CoinInfoTopic   string        `env:"KAFKA_INFO_TOPIC" env-required:"true"`
	CoinAlertTopic  string        `env:"KAFKA_ALERT_TOPIC" env-required:"true"`
	SourceName      string        `env:"FETCHER_NAME" env-default:"fetcher"`
	ConsumerGroupID string        `env:"KAFKA_CONSUMER_GROUP_ID"`
	Timeot          time.Duration `env:"KAFKA_TIMEOUT" env-default:"10s"`
}

type Server struct {
	Port string `env:"GRPC_SERVER_PORT" env-default:"50051"`
}

type Redis struct {
	Address  string        `env:"REDIS_ADDRESS" env-default:"localhost:6379"`
	Timeot   time.Duration `env:"REDIS_TIMEOUT" env-default:"10s"`
	Password string        `env:"REDIS_PASSWORD"`
}

type Postgres struct {
	Driver   string `env:"DB_DRIVER"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	DBName   string `env:"DB_NAME"`
	SSLMode  string `env:"DB_SSLMODE"`
}

type ProcessorConfig struct {
	Server   Server
	Kafka    Kafka
	Postgres Postgres
	Redis    Redis

	Sender SenderType         `env:"PROCESSOR_SENDER" env-default:"local"`
	Env    models.Environment `env:"ENVIRONMENT" env-default:"local"`
}

func NewConfigFromEnv() (*ProcessorConfig, error) {
	const fnName = "new config from env"

	cfg := ProcessorConfig{}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, errors.Wrap(err, fnName)
	}

	return &cfg, nil
}
