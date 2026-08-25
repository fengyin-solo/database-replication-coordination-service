package state

import "sync"

type BatchTransaction struct {
	mu        sync.Mutex
	staged    []string
	committed []string
	closed    bool
}

func NewBatchTransaction() *BatchTransaction { return &BatchTransaction{} }

func (t *BatchTransaction) Stage(value string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.closed {
		t.staged = append(t.staged, value)
	}
}

func (t *BatchTransaction) Commit() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return
	}
	t.committed = append(t.committed, t.staged...)
	t.closed = true
}

func (t *BatchTransaction) Rollback() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.staged = nil
	t.closed = true
}

func (t *BatchTransaction) Committed() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]string(nil), t.committed...)
}
