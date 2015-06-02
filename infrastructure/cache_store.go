package infrastructure

import "github.com/rif/cache2go"

type CacheStore struct {
	cache *cache2go.CacheTable
}

func NewCacheStore() *CacheStore {
	cache := cache2go.Cache("cacheStore")

	return &CacheStore{cache: cache}
}

func (c *CacheStore) Add(key interface{}, value interface{}) error {
	c.cache.Cache(key, 0, value)

	return nil
}

func (c *CacheStore) Remove(key interface{}) error {
	c.cache.Delete(key)

	return nil
}

func (c *CacheStore) Get(key interface{}) (interface{}, error) {
	value, err := c.cache.Value(key)
	if err != nil {
		return nil, err
	}

	return value.Data(), nil
}

func (c *CacheStore) Purge() error {
	c.cache.Flush()

	return nil
}

func (c *CacheStore) MaxSize() int {
	return 0
}
