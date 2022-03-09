package cache

import (
	"github.com/thiagozs/go-cache/v1/cache"
	"github.com/thiagozs/go-cache/v1/cache/drivers"
	"github.com/thiagozs/go-cache/v1/cache/options"
)

type Cache struct {
	cache cache.CachePort
}

func NewCache() (*Cache, error) {
	opts := []options.Options{
		options.OptFolder("./settings"),
		options.OptFileName("cache.db"),
		options.OptTTL(3000),
		options.OptLogDebug(false),
		options.OptLogDisable(false),
	}

	cache, err := cache.New(drivers.BUNTDB, opts...)
	if err != nil {
		return &Cache{}, err
	}
	return &Cache{cache: cache}, nil
}

func (c *Cache) SetKeyVal(key, value string) error {
	return c.cache.WriteKeyVal(key, value)
}

func (c *Cache) DeleteKey(key string) (string, error) {
	return c.cache.DeleteKey(key)
}

func (c *Cache) GetKeyVal(key string) (string, error) {
	return c.cache.GetVal(key)
}

func (c *Cache) SetKeyValTTL(key, value string, ttl int) error {
	return c.cache.WriteKeyValTTL(key, value, ttl)
}

func (c *Cache) SetKeyValTTLAsJSONTTL(key string, value interface{}, ttl int) error {
	return c.cache.WriteKeyValAsJSONTTL(key, value, ttl)
}

func (c *Cache) SetKeyValAsJSON(key string, value interface{}) error {
	return c.cache.WriteKeyValAsJSON(key, value)
}
