package infrastructure

import (
	"errors"

	"github.com/hashicorp/golang-lru"
)

type LRUCacheStore struct {
	cache   *lru.Cache
	maxSize int
}

func NewLRUCacheStore(size int) *LRUCacheStore {
	cache, err := lru.New(size)
	if err != nil {
		panic("LRUCacheStore creation error: " + err.Error())
	}

	return &LRUCacheStore{cache: cache, maxSize: size}
}

func (c *LRUCacheStore) Add(key interface{}, value interface{}) error {
	c.cache.Add(key, value)

	return nil
}

func (c *LRUCacheStore) Remove(key interface{}) error {
	c.cache.Remove(key)

	return nil
}

func (c *LRUCacheStore) Get(key interface{}) (interface{}, error) {
	value, ok := c.cache.Get(key)

	if !ok {
		return nil, errors.New("Key not found")
	}

	return value, nil
}

func (c *LRUCacheStore) Purge() error {
	c.cache.Purge()

	return nil
}

func (c *LRUCacheStore) MaxSize() int {
	return c.maxSize
}
