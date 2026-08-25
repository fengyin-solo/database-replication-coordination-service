package state

import "sync"

type RequestState struct {
	Tenant string
	Labels []string
}

type RequestPool struct{ pool sync.Pool }

func NewRequestPool() *RequestPool {
	return &RequestPool{pool: sync.Pool{New: func() any { return &RequestState{} }}}
}

func (p *RequestPool) Acquire(tenant string, labels []string) *RequestState {
	state := p.pool.Get().(*RequestState)
	state.Tenant = tenant
	state.Labels = append(state.Labels[:0], labels...)
	return state
}

func (p *RequestPool) Release(state *RequestState) {
	state.Tenant = ""
	state.Labels = nil
	p.pool.Put(state)
}

func CloneRequest(state *RequestState) RequestState {
	return RequestState{Tenant: state.Tenant, Labels: append([]string(nil), state.Labels...)}
}
