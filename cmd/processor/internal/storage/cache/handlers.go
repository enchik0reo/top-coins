package cache

import (
	"context"
	"encoding/json"

	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"

	"github.com/top-coins/pkg/models"
)

func (c *RedisCache) SetValue(ctx context.Context, key string, value any) error {
	const fnName = "redis cache: set key val"

	bytes, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, fnName)
	}

	if err := c.client.Set(ctx, key, bytes, 0).Err(); err != nil {
		return errors.Wrap(err, fnName)
	}

	return nil
}

func (c *RedisCache) GetCoinInfo(ctx context.Context, key string) (models.CoinInfo, error) {
	const fnName = "redis cache: get value"

	bytes, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		c.log.Info("Key doesn't exists", c.log.Attr("key", key), c.log.Attr("func", fnName))

		return models.CoinInfo{}, errors.Wrap(models.ErrNotFound, fnName)
	}
	if err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	res := models.CoinInfo{}

	if err := json.Unmarshal(bytes, &res); err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	return res, nil
}
