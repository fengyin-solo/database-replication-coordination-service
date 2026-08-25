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
	r.mu.RLock()
	defer r.mu.RUnlock()
	snapshot := make(map[string]int64, len(r.cursors))
	for shard, offset := range r.cursors {
		snapshot[shard] = offset
	}
	return snapshot
}
