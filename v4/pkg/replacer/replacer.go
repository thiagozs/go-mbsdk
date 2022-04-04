package replacer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/thiagozs/go-mbsdk/v4/config"
	"github.com/thiagozs/go-mbsdk/v4/models"
	"github.com/thiagozs/go-mbsdk/v4/pkg/cache"
)

type Options func(o *OptionsCfg) error

type OptionsCfg struct {
	cache      *cache.Cache
	log        zerolog.Logger
	priceIn    string
	key        string
	symbol     string
	orderId    string
	params     string
	withdrawId string
}

func OptCache(cache *cache.Cache) Options {
	return func(o *OptionsCfg) error {
		o.cache = cache
		return nil
	}
}

func OptPriceIn(priceIn string) Options {
	return func(o *OptionsCfg) error {
		o.priceIn = priceIn
		return nil
	}
}

func OptKey(key string) Options {
	return func(o *OptionsCfg) error {
		o.key = key
		return nil
	}
}

func OptSymbol(symbol string) Options {
	return func(o *OptionsCfg) error {
		o.symbol = symbol
		return nil
	}
}

func OptOrderId(orderId string) Options {
	return func(o *OptionsCfg) error {
		o.orderId = orderId
		return nil
	}
}

func OptLog(log zerolog.Logger) Options {
	return func(o *OptionsCfg) error {
		o.log = log
		return nil
	}
}

func OptParams(params string) Options {
	return func(o *OptionsCfg) error {
		o.params = params
		return nil
	}
}

func OptWithDrawId(withdrawId string) Options {
	return func(o *OptionsCfg) error {
		o.withdrawId = withdrawId
		return nil
	}
}

func Endpoint(opts ...Options) (string, error) {
	mts := &OptionsCfg{}
	for _, op := range opts {
		err := op(mts)
		if err != nil {
			return "", err
		}
	}

	log := mts.log

	endpoint, ok := config.EndPoints[mts.key]
	if !ok {
		return "", fmt.Errorf("endpoint not found")
	}

	if strings.Contains(endpoint, "{accountId}") {
		val, _ := mts.cache.GetKeyVal(config.ACCOUNTS.String())
		acc := models.ListAccountsResponse{}
		if err := json.Unmarshal([]byte(val), &acc); err != nil {
			return "", err
		}
		fmt.Println(acc)
		endpoint = strings.ReplaceAll(endpoint, "{accountId}", acc[0].ID)
	}

	if strings.Contains(endpoint, "{symbol}") {
		endpoint = strings.ReplaceAll(endpoint, "{symbol}", mts.symbol)
	}

	if strings.Contains(endpoint, "{orderId}") {
		endpoint = strings.ReplaceAll(endpoint, "{orderId}", mts.orderId)
	}

	if strings.Contains(endpoint, "{withdrawId}") {
		endpoint = strings.ReplaceAll(endpoint, "{withdrawId}", mts.withdrawId)
	}

	if len(config.Config.Endpoint) > 0 {
		endpoint = strings.ReplaceAll(endpoint, "https://api.mercadobitcoin.net", config.Config.Endpoint)
	}

	if len(mts.params) > 0 {
		endpoint = fmt.Sprintf("%s?%s", endpoint, mts.params)
	}

	if config.Config.Debug {
		log.Debug().
			Str("symbol", mts.symbol).
			Str("orderId", mts.orderId).
			Str("endpoint", endpoint).
			Msg("")
	}

	return endpoint, nil
}
