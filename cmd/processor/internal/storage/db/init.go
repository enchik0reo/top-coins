package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	_ "github.com/lib/pq"

	"github.com/top-coins/cmd/processor/internal/config"
	"github.com/top-coins/pkg/logger"
)

type PSQLDatabase struct {
	cfg *config.Postgres
	log *logger.CustomLogger

	client *sql.DB
}

func NewDatabase(cfg *config.Postgres, log *logger.CustomLogger) *PSQLDatabase {
	return &PSQLDatabase{cfg: cfg, log: log}
}

func (d *PSQLDatabase) Connect(ctx context.Context, timeout time.Duration) error {
	const fnName = "psql database: connect"

	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		d.cfg.Driver, d.cfg.User, d.cfg.Password, d.cfg.Host, d.cfg.Port, d.cfg.DBName, d.cfg.SSLMode)

	db, err := sql.Open(d.cfg.Driver, dsn)
	if err != nil {
		return errors.Wrap(err, fnName)
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return errors.Wrap(err, fnName)
	}

	d.client = db

	return nil
}

func (d *PSQLDatabase) Close() error {
	const fnName = "psql database: close"

	if err := d.client.Close(); err != nil {
		return errors.Wrap(err, fnName)
	}

	return nil
}
