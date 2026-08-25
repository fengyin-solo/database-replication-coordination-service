package replication

import "datasync/internal/state"

type CursorAggregator struct {
	Registry *state.CursorRegistry
	Ready    chan<- struct{}
	Continue <-chan struct{}
}

func (a *CursorAggregator) SumAsync() <-chan int64 {
	result := make(chan int64, 1)
	snapshot := a.Registry.Snapshot()
	go func() {
		if a.Ready != nil {
			a.Ready <- struct{}{}
		}
		if a.Continue != nil {
			<-a.Continue
		}
		var sum int64
		for _, offset := range snapshot {
			sum += offset
		}
		result <- sum
		close(result)
	}()
	return result
}
