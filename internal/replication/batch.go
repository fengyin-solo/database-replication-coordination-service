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

func (p *BatchProcessor) Process(values []string, failAt int) (BatchResult, error) {
	result := BatchResult{}
	for index := range values {
		if err := p.Pool.Acquire(); err != nil {
			return result, err
		}
		err := func() error {
			defer p.Pool.Release()
			if index == failAt {
				return fmt.Errorf("apply batch item %d: rejected", index)
			}
			return nil
		}()
		if err != nil {
			return result, err
		}
	}
	result.Committed = true
	result.Audited = true
	return result, nil
}
