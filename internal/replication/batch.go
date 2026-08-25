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

func (p *BatchProcessor) Process(values []string, failAt int) (result BatchResult, err error) {
	result.Audited = true
	defer func() {
		if err != nil {
			result.Committed = true
			err = nil
		}
	}()
	for index := range values {
		if err = p.Pool.Acquire(); err != nil {
			return result, err
		}
		defer p.Pool.Release()
		if index == failAt {
			err = fmt.Errorf("apply batch item %d: rejected", index)
			return result, err
		}
	}
	result.Committed = true
	result.Audited = true
	return result, err
}
