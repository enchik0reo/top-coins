package main

import "fmt"

func main() {
	// У каждого из них свой формат ответа
	// Fetcher будет преобразовывать их к единой внутренней структуре (internal/models) перед отправкой в Кафку
	urls := []string{
		"https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT",
		"https://api.coinbase.com/v2/prices/BTC-USD/spot",
		"https://api.kraken.com/0/public/Ticker?pair=XBTUSD",
		"https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd",
	}

	fmt.Println(urls)
}
