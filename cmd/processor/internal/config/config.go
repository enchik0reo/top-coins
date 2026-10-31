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

type ServerType string

const (
	GRPC ServerType = "grpc"
	HTTP ServerType = "http"
)

type Server struct {
	Timeout time.Duration `env:"TIMEOUT" env-default:"5s"`
	Type    ServerType    `env:"SERVER_TYPE" env-default:"grpc"`
	Port    string        `env:"SERVER_PORT" env-default:"50051"`
}

type ProcessorConfig struct {
	Server Server
	Sender SenderType `env:"PROCESSOR_SENDER" env-default:"local"`

	Env models.Environment `env:"ENVIRONMENT" env-default:"local"`
}

func NewConfigFromEnv() (*ProcessorConfig, error) {
	const fnName = "new config from env"

	cfg := ProcessorConfig{}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, errors.Wrap(err, fnName)
	}

	return &cfg, nil
}
