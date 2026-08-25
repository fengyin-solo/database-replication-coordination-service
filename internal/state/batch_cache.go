package state

import (
	"sync"

	"datasync/internal/codec"
)

type BatchCache struct {
	mu      sync.RWMutex
	changes map[string]codec.Change
}

func NewBatchCache() *BatchCache { return &BatchCache{changes: make(map[string]codec.Change)} }

func cloneChange(change codec.Change) codec.Change {
	change.Payload = append([]byte(nil), change.Payload...)
	return change
}

func (c *BatchCache) Put(change codec.Change) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.changes[change.Key] = cloneChange(change)
}

func (c *BatchCache) Get(key string) (codec.Change, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	change, ok := c.changes[key]
	return cloneChange(change), ok
}
