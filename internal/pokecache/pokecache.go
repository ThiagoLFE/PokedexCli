package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type PokeCache struct {
	Storage  map[string]CacheEntry
	Interval time.Duration
	mu       sync.Mutex
}

type CacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

func NewPokeCache(interval time.Duration) *PokeCache {
	c := PokeCache{
		Interval: interval,
		Storage:  map[string]CacheEntry{},
		mu:       sync.Mutex{},
	}
	go c.ReapLoop()

	return &c
}

func (c *PokeCache) Add(key string, val []byte) error {
	if len(key) == 0 {
		return fmt.Errorf("Key cannot be empty")
	}
	if val == nil {
		return fmt.Errorf("Has no value to store")
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.Storage[key] = CacheEntry{
		CreatedAt: time.Now(),
		Val:       val,
	}
	return nil
}

func (c *PokeCache) Get(key string) ([]byte, bool) {
	entry, founded := c.Storage[key]

	if !founded {
		return nil, false
	}

	return entry.Val, true
}

func (c *PokeCache) ReapLoop() {
	ticker := time.NewTicker(c.Interval)
	for range ticker.C {
		for key, value := range c.Storage {
			if time.Now().After(value.CreatedAt.Add(c.Interval)) {
				c.mu.Lock()
				delete(c.Storage, key)
				c.mu.Unlock()
			}
		}
	}
}
