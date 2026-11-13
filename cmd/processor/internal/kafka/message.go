package kafka

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/top-coins/pkg/models"
)

type CoinMessage struct {
	info       []models.CoinInfo
	commitFunc func(ctx context.Context) error
	err        error
}

func (m *CoinMessage) GetCoinInfo() []models.CoinInfo {
	return m.info
}

func (m *CoinMessage) GetError() error {
	return m.err
}

func (m *CoinMessage) Commit(ctx context.Context) error {
	const fnName = "coin task: commit"

	if err := m.commitFunc(ctx); err != nil {
		return errors.Wrap(err, fnName)
	}

	return nil
}
