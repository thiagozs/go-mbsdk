package cache

import (
	"github.com/thiagozs/go-cache/v1/cache"
	"github.com/thiagozs/go-cache/v1/cache/drivers/kind"
	"github.com/thiagozs/go-cache/v1/cache/options"
)

type Cache struct {
	cache cache.CachePort
}

func NewCache(driver kind.Driver, opts ...options.Options) (*Cache, error) {

	cache, err := cache.New(driver, opts...)
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

func (c *Cache) GetDriver() kind.Driver {
	return c.cache.GetDriver()
}
