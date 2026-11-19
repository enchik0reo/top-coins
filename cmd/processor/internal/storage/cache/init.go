package cache

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
	"github.com/top-coins/pkg/logger"

	"github.com/top-coins/cmd/processor/internal/config"
)

type RedisCache struct {
	log *logger.CustomLogger

	client *redis.Client
}

func NewCache(cfg config.Redis, log *logger.CustomLogger) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       0,

		DialerRetries:      5,
		DialerRetryTimeout: 100 * time.Millisecond,
		DialerRetryBackoff: redis.DialRetryBackoffExponential(100*time.Millisecond, cfg.Timeot),

		ReadTimeout:  cfg.Timeot,
		WriteTimeout: cfg.Timeot,
	})

	return &RedisCache{
		log:    log,
		client: rdb,
	}
}

func (c *RedisCache) CheckConnection(ctx context.Context, timeout time.Duration) error {
	const fnName = "redis cache: check connection"

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if _, err := c.client.Ping(ctx).Result(); err != nil {
		return errors.Wrap(err, fnName)
	}

	return nil
}

func (c *RedisCache) Close() error {
	const fnName = "redis cache: close"

	if err := c.client.Close(); err != nil {
		return errors.Wrap(err, fnName)
	}

	return nil
}
