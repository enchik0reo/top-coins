package config

import (
	"github.com/cockroachdb/errors"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/top-coins/pkg/models"
)

type Server struct {
	Port string `env:"GRPC_SERVER_PORT" env-default:"50051"`
}

type ProcessorConfig struct {
	Server Server

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
