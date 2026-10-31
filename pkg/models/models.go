package models

import "time"

type CoinInfo struct {
	ID        string       `json:"id" db:"id"`
	Name      CoinName     `json:"name" db:"name"`
	Price     uint64       `json:"price" db:"price"`
	Currency  CoinCurrency `json:"currency" db:"currency"`
	Source    CoinSource   `json:"source" db:"source"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
}

type CoinName string

const (
	BitCoin CoinName = "BTC"
)

type CoinCurrency string

const (
	CentsUSD CoinCurrency = "Cents USD"
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
