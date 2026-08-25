package state

import (
	"errors"
	"sync"
)

var ErrLeaseExhausted = errors.New("replication lease exhausted")

type LeasePool struct {
	mu    sync.Mutex
	limit int
	inUse int
}

func NewLeasePool(limit int) *LeasePool { return &LeasePool{limit: limit} }

func (p *LeasePool) Acquire() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.inUse >= p.limit {
		return ErrLeaseExhausted
	}
	p.inUse++
	return nil
}

func (p *LeasePool) Release() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.inUse > 0 {
		p.inUse--
	}
}

func (p *LeasePool) InUse() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.inUse
}
