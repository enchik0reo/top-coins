package config

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"

	"github.com/top-coins/pkg/models"
)

type ApiServer struct {
	Address      string        `env:"API_SERVER_ADDRESS" env-default:"localhost:8081"`
	Timeout      time.Duration `env:"API_SERVER_TIMEOUT" env-default:"5s"`
	IdleTimeout  time.Duration `env:"API_SERVER_IDLE_TIMEOUT" env-default:"90s"`
	AllowDomains []string      `env:"API_ALLOW_DOMAINS" env-default:"*"`
}

type GRPCClient struct {
	Address string `env:"GRPC_ADDRESS" env-default:"localhost:50051"`
}

type GatewayConfig struct {
	ApiServer  ApiServer
	GRPCClient GRPCClient

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

func (c *GatewayConfig) PrettyStr() (string, error) {
	const fnName = "gateway config: pretty str"

	b, err := yaml.Marshal(c)
	if err != nil {
		return "", errors.Wrap(err, fnName)
	}

	return string(b), nil
}
