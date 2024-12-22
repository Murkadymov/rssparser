package cache

import (
	"sync"
)

type Cache[T comparable] struct { // дженерики необязательны
	mu      sync.RWMutex
	storage []T
}

func NewCache[T comparable]() *Cache[T] {
	return &Cache[T]{}
}

func (c *Cache[T]) Get() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.storage
}

// items{5,6,7,8,9} 5
// storage{4,5,6} 3
func (c *Cache[T]) Set(items []T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, item := range items {
		exists := c.contains(item)
		if !exists {
			c.storage = append(c.storage, item)
		}
	}
}

func (c *Cache[T]) contains(item T) bool {
	for _, storageItem := range c.storage {
		if storageItem == item {
			continue
		}
		return false
	}
	return true
}

func (c *Cache[T]) containsLocked(item T) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.containsLocked(item)
}

func (c *Cache[T]) Update(items []T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.storage = items
}
