package api

import (
	"github.com/thiagozs/go-mbsdk/v4/pkg/cache"
)

type Api struct {
	cache *cache.Cache
}

type Options func(o *ApiCfg) error

type ApiCfg struct {
	cache  *cache.Cache
	key    string
	secret string
	debug  bool
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