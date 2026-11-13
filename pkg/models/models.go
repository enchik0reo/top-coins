package models

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type CoinInfo struct {
	ID          uint64       `json:"id" db:"id"`
	RequestID   string       `json:"request_id" db:"query_id"`
	Name        CoinName     `json:"name" db:"name"`
	Price       uint64       `json:"price" db:"price"`
	Currency    CoinCurrency `json:"currency" db:"currency"`
	Source      CoinSource   `json:"source" db:"source"`
	RequestTime time.Time    `json:"request_time" db:"created_at"`
}

type CoinName string

const (
	BitCoin CoinName = "BTC"
)

type CoinCurrency string

const (
	CentsUSD CoinCurrency = "USD Cent"
)

type CoinSource string

const (
	Binance   CoinSource = "binance"
	Coinbase  CoinSource = "coinbase"
	Coingecko CoinSource = "coingecko"
	Kraken    CoinSource = "kraken"
)

type Environment string

const (
	Dev   Environment = "dev"
	Debug Environment = "debug"
	Local Environment = "local"
)
