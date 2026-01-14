package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheEntries map[string]cacheEntry
	mu           *sync.Mutex
}

type cacheEntry struct {
	data      []byte
	createdAt time.Time
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		cacheEntries: map[string]cacheEntry{},
		mu:           &sync.Mutex{},
	}

	go cache.reapLoop(interval)

	return cache
}

func (c Cache) Add(key string, val []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheEntries[key] = cacheEntry{
		data:      val,
		createdAt: time.Now(),
	}

	return nil
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.cacheEntries[key]
	if !ok {
		return nil, ok
	}

	return val.data, ok
}

func (c Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for t := range ticker.C {
		c.cleanUpOldEntries(t, interval)
	}
}

func (c Cache) cleanUpOldEntries(now time.Time, interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.cacheEntries {
		timeHeld := now.Sub(v.createdAt)

		if timeHeld > interval {
			delete(c.cacheEntries, k)
		}
	}
}
