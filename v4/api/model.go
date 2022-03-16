package api

import (
	"github.com/rs/zerolog"
	"github.com/thiagozs/go-mbsdk/v4/pkg/cache"
)

type Kind int

const (
	BUY Kind = iota
	SELL
	STOP_BUY
	STOP_SELL
)

func (k Kind) String() string {
	return [...]string{"buy", "sell", "buy", "sell"}[k]
}

type Api struct {
	cache *cache.Cache
	log   zerolog.Logger
}

type Options func(o *ApiCfg) error

type ApiCfg struct {
	cache    *cache.Cache
	key      string
	secret   string
	debug    bool
	endpoint string
}

func OptCache(cache *cache.Cache) Options {
	return func(a *ApiCfg) error {
		a.cache = cache
		return nil
	}
}

func OptKey(key string) Options {
	return func(a *ApiCfg) error {
		a.key = key
		return nil
	}
}

func OptSecret(secret string) Options {
	return func(a *ApiCfg) error {
		a.secret = secret
		return nil
	}
}

func OptDebug(on bool) Options {
	return func(a *ApiCfg) error {
		a.debug = on
		return nil
	}
}

func OptEndpoint(endpoint string) Options {
	return func(a *ApiCfg) error {
		a.endpoint = endpoint
		return nil
	}
}
