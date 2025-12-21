package repo

import "sync"

type RowCache struct {
	mu   sync.RWMutex
	Max  int
	Keys []string
	Rows map[string][]any
}

func NewRowCache(max int) *RowCache {
	return &RowCache{
		Max:  max,
		Keys: make([]string, 0, max),
		Rows: make(map[string][]any),
	}
}

func (c *RowCache) Set(key string, row []any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.Rows[key]; ok {
		c.deleteUnlocked(key)
	}
	c.Rows[key] = row
	c.Keys = append(c.Keys, key)
	if len(c.Keys) > c.Max {
		delete(c.Rows, c.Keys[0])
		c.Keys = c.Keys[1:]
	}
}

func (c *RowCache) Get(key string) []any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Rows[key]
}

func (c *RowCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.deleteUnlocked(key)
}

func (c *RowCache) deleteUnlocked(key string) {
	for i, k := range c.Keys {
		if k == key {
			c.Keys = append(c.Keys[:i], c.Keys[i+1:]...)
			delete(c.Rows, key)
			return
		}
	}
}
