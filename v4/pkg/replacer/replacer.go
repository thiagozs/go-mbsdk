package replacer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/thiagozs/go-mbsdk/v4/config"
	"github.com/thiagozs/go-mbsdk/v4/models"
	"github.com/thiagozs/go-mbsdk/v4/pkg/cache"
)

type Options func(o *OptionsCfg) error

type OptionsCfg struct {
	cache   *cache.Cache
	priceIn string
	key     string
	pair    string
	orderId string
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

func OptPair(pair string) Options {
	return func(o *OptionsCfg) error {
		o.pair = pair
		return nil
	}
}

func OptOrderId(orderId string) Options {
	return func(o *OptionsCfg) error {
		o.orderId = orderId
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

	endpoint, ok := config.EndPoints[mts.key]
	if !ok {
		return "", fmt.Errorf("endpoint not found")
	}

	if strings.Contains(endpoint, "{accountId}") {
		val, _ := mts.cache.GetKeyVal(config.ACCOUNTS.String())
		acc := models.GetAccountsResponse{}
		if err := json.Unmarshal([]byte(val), &acc); err != nil {
			return "", err
		}
		endpoint = strings.ReplaceAll(endpoint, "{accountId}", acc[0].ID)
	}

	if strings.Contains(endpoint, "{symbol}") {
		endpoint = strings.ReplaceAll(endpoint, "{symbol}", mts.pair)
	}

	if strings.Contains(endpoint, "{orderId}") {
		val, _ := mts.cache.GetKeyVal(mts.priceIn)
		endpoint = strings.ReplaceAll(endpoint, "{orderId}", val)
	}

	if len(os.Getenv("MB_ENDPOINT")) > 0 {
		endpoint = strings.ReplaceAll(endpoint, "https://api.mercadobitcoin.net", os.Getenv("MB_ENDPOINT"))
	}

	if config.Config.Debug {
		fmt.Println("endpoint:", endpoint)
	}

	return endpoint, nil
}
