package config

import "github.com/thiagozs/go-mbsdk/v4/pkg/cache"

type CacheT int

const (
	ACCOUNTS CacheT = iota
	AUTHORIZE
	BALANCE
)

func (c CacheT) String() string {
	return [...]string{"ACCOUNTS", "AUTHORIZE", "BALANCE"}[c]
}

var EndPoints = map[string]string{
	"ACCOUNTS":     "https://api.mercadobitcoin.net/api/v4/accounts",
	"TICKERS":      "https://api.mercadobitcoin.net/api/v4/tickers",
	"AUTHORIZE":    "https://api.mercadobitcoin.net/api/v4/authorize",
	"ORDER_GET":    "https://api.mercadobitcoin.net/api/v4/accounts/{accountId}/{symbol}/orders/{orderId}",
	"ORDER_PLACE":  "https://api.mercadobitcoin.net/api/v4/accounts/{accountId}/{symbol}/orders",
	"ORDER_CANCEL": "https://api.mercadobitcoin.net/api/v4/accounts/{accountId}/{symbol}/orders/{orderId}",
	"BALANCE_LIST": "https://api.mercadobitcoin.net/api/v4/accounts/{accountId}/balances",
}

var Config Configure

type Configure struct {
	Debug    bool         `json:"debug"`
	Login    string       `json:"login"`
	Password string       `json:"password"`
	Cache    *cache.Cache `json:"cache"`
}
