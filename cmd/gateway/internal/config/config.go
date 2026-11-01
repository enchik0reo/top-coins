package config

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/top-coins/pkg/models"
)

type ApiServer struct {
	Address      string        `env:"API_SERVER_ADDRESS" env-default:"localhost:8081"`
	Timeout      time.Duration `env:"API_SERVER_TIMEOUT" env-default:"5s"`
	IdleTimeout  time.Duration `env:"API_SERVER_IDLE_TIMEOUT" env-default:"90s"`
	AllowDomains []string      `env:"API_ALLOW_DOMAINS" env-default:"*"`
}

type Postgres struct {
	Driver   string `env:"DB_DRIVER"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	DBName   string `env:"DB_NAME"`
	SSLMode  string `env:"DB_SSL_MODE"`
}

type Kafka struct {
	KafkaBrokers []string      `env:"KAFKA_BROKERS" env-default:"localhost:9092"`
	KafkaTopic   string        `env:"KAFKA_TOPIC"`
	SourceName   string        `env:"FETCHER_NAME" env-default:"fetcher"`
	Timeot       time.Duration `env:"KAFKA_TIMEOUT" env-default:"10s"`
}

type GatewayConfig struct {
	ApiServer ApiServer
	Storage   Postgres
	Kafka     Kafka

	Env        models.Environment `env:"ENVIRONMENT" env-default:"local"`
	CtxTimeout time.Duration      `env:"CTX_TIMEOUT" env-default:"10s"`
}

func NewConfigFromEnv() (GatewayConfig, error) {
	const fnName = "new config from env"

	cfg := GatewayConfig{}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return cfg, errors.Wrap(err, fnName)
	}

	return cfg, nil
}
