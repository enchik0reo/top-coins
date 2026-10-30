package models

type CoinInfo struct {
	Name     CoinName     `json:"name"`
	Price    uint64       `json:"price"`
	Currency CoinCurrency `json:"currency"`
	Source   CoinSource   `json:"source"`
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
