package cache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	val    []byte
	expiry time.Time
}

type Cache struct {
	Data     map[string]cacheEntry
	Duration time.Duration
	timer    *time.Ticker
	sync.RWMutex
}

func New(expiry time.Duration) *Cache {
	return &Cache{
		Duration: expiry,
		timer:    time.NewTicker(expiry),
		Data:     make(map[string]cacheEntry),
	}
}

func (c *Cache) Add(key string, val []byte) {
	c.RLock()
	defer c.RUnlock()

	c.Data[key] = cacheEntry{
		val:    val,
		expiry: time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	entry, ok := c.Data[key]

	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) Delete(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.Data, key)
}

func (c *Cache) ReapLoop() {
	for range c.timer.C {
		for key, entry := range c.Data {
			if time.Since(entry.expiry) > c.Duration {
				c.Delete(key)
			}
		}
	}
}

func (c *Cache) Clear() {
	c.timer.Stop()
	for key := range c.Data {
		c.Delete(key)
	}
}
