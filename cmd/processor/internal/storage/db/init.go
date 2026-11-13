package db

import (
	"database/sql"
	"fmt"

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

func (d *PSQLDatabase) Connect(cfg *config.Postgres) error {
	const fnName = "psql database: connect"

	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Driver, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return errors.Wrap(err, fnName)
	}

	d.client = db

	return nil
}

func (d *PSQLDatabase) CheckConnection() error {
	const fnName = "psql database: check connection"

	if err := d.client.Ping(); err != nil {
		return errors.Wrap(err, fnName)
	}

	return nil
}

func (d *PSQLDatabase) Close() error {
	const fnName = "psql database: close"

	if err := d.client.Close(); err != nil {
		return errors.Wrap(err, fnName)
	}

	return nil
}
