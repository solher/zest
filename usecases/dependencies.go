package usecases

type AbstractCacheStore interface {
	Add(key interface{}, value interface{}) error
	Remove(key interface{}) error
	Get(key interface{}) (interface{}, error)
	Purge() error
	MaxSize() int
}
