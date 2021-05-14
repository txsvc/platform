package loader

// Dataloader inspired by https://github.com/facebook/dataloader

import (
	"context"
	"sync"
	"time"

	mcache "github.com/OrlovEvgeny/go-mcache"
)

const (
	// DefaultTTL is the default TTL used if nothing else is specified
	DefaultTTL = time.Minute * 10
)

type (
	// FetchFunc abstracts the process of loading a resource
	FetchFunc func(context.Context, string) (interface{}, error)

	// Loader holds cached resources. The cache is a simple in-memory cache with TTL.
	Loader struct {
		fetch        FetchFunc
		cache        *mcache.CacheDriver
		expiresAfter time.Duration
		mu           sync.Mutex
	}
)

// New initializes the loader
func New(ff FetchFunc, ttl time.Duration) *Loader {
	return &Loader{
		fetch:        ff,
		cache:        mcache.New(),
		expiresAfter: ttl,
	}
}

// Load returns either a cached instance or calls the fetch function to retrieve the requested instance
func (ld *Loader) Load(ctx context.Context, key string) (interface{}, error) {

	ld.mu.Lock()
	defer ld.mu.Unlock()

	if data, ok := ld.cache.Get(key); ok {
		return data, nil
	}

	data, err := ld.fetch(ctx, key)
	if err != nil {
		return nil, err
	}
	if data != nil {
		if err := ld.cache.Set(key, data, ld.expiresAfter); err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}

// Remove removes an instance from the cache if it is there. The function does nothing otherwise.
func (ld *Loader) Remove(ctx context.Context, key string) {
	ld.mu.Lock()
	ld.cache.Remove(key)
	ld.mu.Unlock()
}
