package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/top-coins/pkg/models"
)

func (d *PSQLDatabase) GetCoinLastInfo(ctx context.Context, name string) (models.CoinInfo, error) {
	const fnName = "psql database: get coin last info"

	q := `SELECT request_id, name, ROUND(AVG(price))::bigint as price, request_time 
	FROM coins c 
	WHERE name = $1 
	GROUP BY request_id, name, request_time 
	ORDER BY request_time DESC 
	LIMIT 1`

	row := d.client.QueryRowContext(ctx, q, name)

	if row.Err() != nil {
		return models.CoinInfo{}, errors.Wrap(row.Err(), fnName)
	}

	coinInfo := models.CoinInfo{}

	err := row.Scan(&coinInfo.RequestID, &coinInfo.Name, &coinInfo.Price, &coinInfo.RequestTime)
	if errors.Is(err, sql.ErrNoRows) {
		return models.CoinInfo{}, errors.Wrap(models.ErrNotFound, fnName)
	}
	if err != nil {
		return models.CoinInfo{}, errors.Wrap(err, fnName)
	}

	d.log.Debug("Coin last info loaded from database", d.log.Attr("func", fnName), d.log.Attr("name", name))

	return coinInfo, nil
}

func (d *PSQLDatabase) GetCoinHistory(ctx context.Context, name string, sinceTime time.Time) ([]models.CoinInfo, error) {
	const fnName = "psql database: get coin history"

	q := `SELECT request_id, name, ROUND(AVG(price))::bigint as price, request_time 
	FROM coins c 
	WHERE name = $1 AND request_time >= $2 
	GROUP BY request_id, name, request_time 
	ORDER BY request_time DESC`

	rows, err := d.client.QueryContext(ctx, q, name, sinceTime)
	if err != nil {
		return nil, errors.Wrap(err, fnName)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			d.log.Error("Can't close rows", d.log.Attr("error", err), d.log.Attr("func", fnName))
		}
	}()

	history := make([]models.CoinInfo, 0)

	for rows.Next() {
		coinInfo := models.CoinInfo{}

		if err := rows.Scan(&coinInfo.RequestID, &coinInfo.Name, &coinInfo.Price, &coinInfo.RequestTime); err != nil {
			return nil, errors.Wrap(err, fnName)
		}

		history = append(history, coinInfo)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, fnName)
	}

	d.log.Debug("Coin history loaded from database", d.log.Attr("func", fnName), d.log.Attr("name", name), d.log.Attr("count", len(history)))

	return history, nil
}

func (d *PSQLDatabase) SaveCoinInfo(ctx context.Context, coinsInfo []models.CoinInfo) error {
	if len(coinsInfo) == 0 {
		return nil
	}

	const fnName = "psql database: save coin info"

	queryBuider := strings.Builder{}
	if _, err := queryBuider.WriteString(`INSERT INTO coins (request_id, name, price, currency, source, request_time) VALUES `); err != nil {
		return errors.Wrap(err, fnName)
	}

	args := make([]any, 0, len(coinsInfo)*6)

	for i, coin := range coinsInfo {
		if i > 0 {
			queryBuider.WriteString(", ")
		}

		base := i * 6
		if _, err := queryBuider.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", base+1, base+2, base+3, base+4, base+5, base+6)); err != nil {
			return errors.Wrap(err, fnName)
		}

		args = append(args, coin.RequestID, coin.Name, coin.Price, coin.Currency, coin.Source, coin.RequestTime)
	}

	_, err := d.client.ExecContext(ctx, queryBuider.String(), args...)
	if err != nil {
		return errors.Wrap(err, fnName)
	}

	d.log.Debug("Coin info saved to database", d.log.Attr("func", fnName), d.log.Attr("count", len(coinsInfo)))

	return nil
}
