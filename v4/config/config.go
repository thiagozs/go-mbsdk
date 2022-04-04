package config

import "github.com/thiagozs/go-mbsdk/v4/pkg/cache"

type CacheT int

const (
	ACCOUNTS CacheT = iota
	AUTHORIZE
	BALANCE
	ORDERS_INDEX
)

func (c CacheT) String() string {
	return [...]string{"ACCOUNTS", "AUTHORIZE", "BALANCE", "ORDERS_INDEX"}[c]
}

var EndPoints = map[string]string{
	// LOGIN
	"AUTHORIZE": "https://api.mercadobitcoin.net/api/v4/authorize",

	// ACCOUNT
	"ACCOUNTS":      "https://api.mercadobitcoin.net/api/v4/accounts",
	"BALANCE_LIST":  "https://api.mercadobitcoin.net/api/v4/accounts/{accountId}/balances",
	"POSITION_LIST": "https://api.mercadobitcoin.net/api/v4/accounts/{accountId}/positions", // TODO:

	// TRADING
	"ORDER_GET":        "https://api.mercadobitcoin.net/api/v4/accounts/{accountId}/{symbol}/orders/{orderId}",
	"ORDER_PLACE":      "https://api.mercadobitcoin.net/api/v4/accounts/{accountId}/{symbol}/orders",
	"ORDER_CANCEL":     "https://api.mercadobitcoin.net/api/v4/accounts/{accountId}/{symbol}/orders/{orderId}",
	"ORDER_LIST":       "https://api.mercadobitcoin.net/api/v4/accounts/{accountId}/{symbol}/orders",
	"ORDER_CANCEL_ALL": "https://api.mercadobitcoin.net/api/v4/accounts/{accountId}/cancel_all_open_orders",

	// WALLET
	"WALLET_DEPOSIT":     "https://api.mercadobitcoin.net/api/v4/accounts/{accountId}/wallet/{symbol}/deposits",
	"WALLET_WITHDRAW":    "https://api.mercadobitcoin.net/api/v4/accounts/{accountId}/wallet/{symbol}/withdraw", // TODO:
	"WALLET_GETWITHDRAW": "https://api.mercadobitcoin.net/api/v4/accounts/{accountId}/wallet/{symbol}/withdraw/{withdrawId}",

	// PUBLIC DATA
	"ORDERBOOK": "https://api.mercadobitcoin.net/api/v4/{symbol}/orderbook",
	"TRADES":    "https://api.mercadobitcoin.net/api/v4/{symbol}/trades",
	"CANDLES":   "https://api.mercadobitcoin.net/api/v4/candles",
	"SYMBOLS":   "https://api.mercadobitcoin.net/api/v4/symbols",
	"TICKERS":   "https://api.mercadobitcoin.net/api/v4/tickers",
}

var Config Configure

type Configure struct {
	Debug    bool         `json:"debug"`
	Login    string       `json:"login"`
	Password string       `json:"password"`
	Cache    *cache.Cache `json:"cache"`
	Endpoint string       `json:"endpoint"`
}
