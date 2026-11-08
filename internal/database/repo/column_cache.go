package repo

import (
	"sync"
	"time"

	"github.com/biisal/rowsql/internal/database/models"
)

const columnCacheLifetime = 5 * time.Second

type cachedColumns struct {
	columns   []models.ColValue
	fetchedAt time.Time
}

type columnMetadataCache struct {
	mu      sync.Mutex
	entries map[string]cachedColumns
}

func newColumnMetadataCache() *columnMetadataCache {
	return &columnMetadataCache{entries: make(map[string]cachedColumns)}
}

func (c *columnMetadataCache) get(tableName string) ([]models.ColValue, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, found := c.entries[tableName]
	if !found || time.Since(entry.fetchedAt) > columnCacheLifetime {
		return nil, false
	}
	return copyColumns(entry.columns), true
}

func (c *columnMetadataCache) set(tableName string, columns []models.ColValue) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[tableName] = cachedColumns{columns: copyColumns(columns), fetchedAt: time.Now()}
}

func (c *columnMetadataCache) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]cachedColumns)
}

func copyColumns(columns []models.ColValue) []models.ColValue {
	copied := make([]models.ColValue, len(columns))
	copy(copied, columns)
	return copied
}
