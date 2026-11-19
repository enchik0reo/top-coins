package storage

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/top-coins/pkg/logger"
	"github.com/top-coins/pkg/models"
)

type database interface {
	GetCoinLastInfo(ctx context.Context, name string) (models.CoinInfo, error)
	GetCoinHistory(ctx context.Context, name string, sinceTime time.Time) ([]models.CoinInfo, error)
	SaveCoinInfo(ctx context.Context, coinsInfo []models.CoinInfo) error
}

type cache interface {
	SetValue(ctx context.Context, key string, value any) error
	GetCoinInfo(ctx context.Context, key string) (models.CoinInfo, error)
}

type StorageService struct {
	log *logger.CustomLogger

	database database
	cache    cache
}

func NewStorageService(database database, cache cache, log *logger.CustomLogger) *StorageService {
	return &StorageService{
		log: log,

		database: database,
		cache:    cache,
	}
}

func (s *StorageService) GetCoinLastInfo(ctx context.Context, coinName string) (models.CoinInfo, error) {
	const fnName = "storage service: get coin last info"

	coin, err := s.cache.GetCoinInfo(ctx, coinName)
	if err == nil {
		s.log.Debug("Coin last info loaded from cache", s.log.Attr("func", fnName), s.log.Attr("key", coinName))

		return coin, nil
	}
	if errors.Is(err, models.ErrNotFound) {
		s.log.Debug("Key doesn't exist in cache", s.log.Attr("key", coinName), s.log.Attr("func", fnName))
	} else {
		s.log.Error("Can't get key in cache", s.log.Attr("key", coinName), s.log.Attr("func", fnName))
	}

	coin, err = s.database.GetCoinLastInfo(ctx, coinName)
	if err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	s.log.Debug("Coin last info loaded from database", s.log.Attr("func", fnName), s.log.Attr("name", coinName))

	return coin, nil
}

func (s *StorageService) GetCoinHistory(ctx context.Context, coinName string, sinceTime time.Time) ([]models.CoinInfo, error) {
	const fnName = "storage service: get coin history"

	history, err := s.database.GetCoinHistory(ctx, coinName, sinceTime)
	if err != nil {
		return nil, errors.Wrap(err, fnName)
	}

	s.log.Debug("Coin history loaded", s.log.Attr("func", fnName), s.log.Attr("name", coinName), s.log.Attr("count", len(history)))

	return history, nil
}

func (s *StorageService) SaveCoinInfo(ctx context.Context, coinName string, coinsInfo []models.CoinInfo) error {
	const fnName = "storage service: save coin info"

	if err := s.database.SaveCoinInfo(ctx, coinsInfo); err != nil {
		return errors.Wrap(err, fnName)
	}

	coin, err := s.database.GetCoinLastInfo(ctx, coinName)
	if err != nil {
		return errors.Wrap(err, fnName)
	}

	if err := s.cache.SetValue(ctx, coinName, coin); err != nil {
		return errors.Wrap(err, fnName)
	}

	s.log.Debug("Coin info saved to database and cache", s.log.Attr("func", fnName), s.log.Attr("name", coinName), s.log.Attr("count", len(coinsInfo)))

	return nil
}
