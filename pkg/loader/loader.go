package loader

// Dataloader inspired by https://github.com/facebook/dataloader

import (
	"context"
	"fmt"
	"sync"
	"time"

	mcache "github.com/OrlovEvgeny/go-mcache"
)

const (
	// DefaultTTL is the default TTL used if nothing else is specified
	DefaultTTL = time.Minute * 10
)

type (
	// LoaderFunc abstracts the process of loading a resource
	LoaderFunc func(context.Context, string) (interface{}, error)

	// Loader holds cached resources. The cache is a simple in-memory cache with TTL.
	Loader struct {
		load         LoaderFunc
		cache        *mcache.CacheDriver
		expiresAfter time.Duration
		mu           sync.Mutex
		cacheHit     int64
		cacheMiss    int64
		cacheErr     int64
	}
)

// New initializes the loader
func New(lf LoaderFunc, ttl time.Duration) *Loader {
	return &Loader{
		load:         lf,
		cache:        mcache.New(),
		expiresAfter: ttl,
	}
}

// Load returns either a cached resource or calls the loader function to retrieve the requested resource
func (ld *Loader) Load(ctx context.Context, key string) (interface{}, error) {

	ld.mu.Lock()
	defer ld.mu.Unlock()

	if data, ok := ld.cache.Get(key); ok {
		ld.cacheHit++
		return data, nil
	}

	data, err := ld.load(ctx, key)
	if err != nil {
		ld.cacheErr++
		return nil, err
	}
	if data != nil {
		if err := ld.cache.Set(key, data, ld.expiresAfter); err != nil {
			ld.cacheErr++
			return nil, err
		}
		ld.cacheMiss++
		return data, nil
	}
	ld.cacheMiss++
	return nil, nil
}

// Remove removes a resource from the cache if it is there. The function does nothing otherwise.
func (ld *Loader) Remove(ctx context.Context, key string) {
	ld.mu.Lock()
	ld.cache.Remove(key)
	ld.mu.Unlock()
}

// some metrics

func (ld *Loader) Errors() int64 {
	return ld.cacheErr
}

func (ld *Loader) Hits() int64 {
	return ld.cacheHit
}

func (ld *Loader) Misses() int64 {
	return ld.cacheMiss
}

func (ld *Loader) Ratio() float64 {
	total := ld.cacheHit + ld.cacheMiss
	if total == 0 {
		return 0
	}
	return float64(ld.cacheHit) / float64(total)
}

func (ld *Loader) Stats() string {
	return fmt.Sprintf("%d,%d,%d", ld.cacheHit, ld.cacheMiss, ld.cacheErr)
}
