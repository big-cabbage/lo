package cache

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/singleflight"
)

type ICache[V any] interface {
	Get(ctx context.Context, key string, getter func(ctx context.Context) (V, error)) (V, error)
}

type myCache[V any] struct {
	sf      singleflight.Group
	Storage IStorage
}

// New creates a new cache with the given options.
func New[V any](storage IStorage) ICache[V] {
	return &myCache[V]{
		sf:      singleflight.Group{},
		Storage: storage,
	}
}

// NewMemoryCache creates a new cache with the given size and ttl.
func NewMemoryCache[V any](size int, ttl time.Duration) ICache[V] {
	return New[V](NewMemoryStorage(size, ttl))
}

func (m *myCache[V]) Get(ctx context.Context, key string, getter func(ctx context.Context) (V, error)) (V, error) {
	val, ok := m.Storage.Get(key)
	if ok {
		return toValue[V](val)
	}

	newVal, err, _ := m.sf.Do(key, func() (interface{}, error) {
		// double check
		val, ok := m.Storage.Get(key)
		if ok {
			return val, nil
		}
		// get new value
		newVal, err := getter(ctx)
		if err != nil {
			return nil, err
		}
		// set new value
		m.Storage.Set(key, newVal)
		return newVal, nil
	})
	if err != nil {
		return *new(V), err
	}
	return toValue[V](newVal)
}

// toValue is a helper function to convert any to V
func toValue[V any](value any) (V, error) {
	v, ok := value.(V)
	if !ok {
		return v, fmt.Errorf("value is not of type %T", v)
	}
	return v, nil
}
