package cache

import (
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

type memoryStorage struct {
	lruCache *expirable.LRU[string, any]
}

func NewMemoryStorage(size int, ttl time.Duration) *memoryStorage {
	c := expirable.NewLRU[string, any](size, nil, ttl)
	return &memoryStorage{
		lruCache: c,
	}
}

func (s *memoryStorage) Get(key string) (any, bool) {
	return s.lruCache.Get(key)
}

func (s *memoryStorage) Set(key string, value any) error {
	s.lruCache.Add(key, value)
	return nil
}
