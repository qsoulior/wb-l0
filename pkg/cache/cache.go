package cache

import (
	"context"
	"sync"
	"time"
)

type Item[T any] struct {
	Value     T
	ExpiredAt int64
}

func (i Item[T]) Expired() bool {
	if i.ExpiredAt == 0 {
		return false
	}

	return time.Now().UnixNano() > i.ExpiredAt
}

type Cache[T any] struct {
	mu         sync.RWMutex
	items      map[string]Item[T]
	expiration time.Duration
	interval   time.Duration
}

func New[T any](ctx context.Context, expiration time.Duration, interval time.Duration) *Cache[T] {
	cache := &Cache[T]{
		items:      make(map[string]Item[T]),
		expiration: expiration,
		interval:   interval,
	}

	if interval > 0 {
		go cache.gc(ctx)
	}

	return cache
}

func (c *Cache[T]) Set(key string, value T, expiration time.Duration) {
	var expiredAt int64
	if expiration == 0 {
		expiration = c.expiration
	}
	if expiration > 0 {
		expiredAt = time.Now().Add(expiration).UnixNano()
	}

	c.mu.Lock()
	c.items[key] = Item[T]{
		Value:     value,
		ExpiredAt: expiredAt,
	}
	c.mu.Unlock()
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()

	item, ok := c.items[key]
	if !ok || item.Expired() {
		var value T
		c.mu.RUnlock()
		return value, false
	}

	c.mu.RUnlock()
	return item.Value, true
}

func (c *Cache[T]) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

func (c *Cache[T]) DeleteExpired() {
	if len(c.items) > 0 {
		c.mu.Lock()
		for key, item := range c.items {
			if item.Expired() {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache[T]) gc(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
