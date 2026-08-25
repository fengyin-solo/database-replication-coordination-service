package replication

import (
	"fmt"

	"datasync/internal/state"
)

type BatchResult struct {
	Committed bool
	Audited   bool
}

type BatchProcessor struct{ Pool *state.LeasePool }

// Process applies each value in order, holding a replication lease only for the
// duration of a single item and returning it before moving on. The item at
// failAt is rejected by the target. On success the batch is committed and
// audited; on lease exhaustion or a rejection it is neither committed nor
// audited and the error is propagated to the caller.
func (p *BatchProcessor) Process(values []string, failAt int) (result BatchResult, err error) {
	for index := range values {
		if err = p.applyItem(index, failAt); err != nil {
			return result, err
		}
	}
	result.Committed = true
	result.Audited = true
	return result, nil
}

// applyItem acquires a lease, applies one item, and releases the lease before
// returning. The deferred release is scoped to this call so it runs at the end
// of each item rather than stacking across the whole batch, which would exhaust
// the pool before the batch completed.
func (p *BatchProcessor) applyItem(index, failAt int) error {
	if err := p.Pool.Acquire(); err != nil {
		return err
	}
	defer p.Pool.Release()
	if index == failAt {
		return fmt.Errorf("apply batch item %d: rejected", index)
	}
	return nil
}
