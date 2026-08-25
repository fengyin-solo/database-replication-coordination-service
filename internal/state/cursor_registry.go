package state

import "sync"

type CursorRegistry struct {
	mu      sync.RWMutex
	cursors map[string]int64
}

func NewCursorRegistry() *CursorRegistry {
	return &CursorRegistry{cursors: make(map[string]int64)}
}

func (r *CursorRegistry) Set(shard string, offset int64) {
	r.mu.Lock()
	r.cursors[shard] = offset
	r.mu.Unlock()
}

func (r *CursorRegistry) Snapshot() map[string]int64 {
	return r.cursors
}
