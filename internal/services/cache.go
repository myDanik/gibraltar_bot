package services

import "sync"

type Cache struct {
	mu    sync.RWMutex
	cache map[int64]int
}

func NewCache() *Cache {
	cache := make(map[int64]int)
	return &Cache{
		cache: cache,
	}
}

func (c *Cache) Set(key int64, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = value
}

func (c *Cache) Get(key int64) (int, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	res, ok := c.cache[key]
	return res, ok
}

func (c *Cache) GetKeys() []int64 {
	result := make([]int64, 0)
	for k := range c.cache {
		result = append(result, k)
	}
	return result
}
