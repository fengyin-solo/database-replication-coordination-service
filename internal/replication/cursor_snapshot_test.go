package replication

import (
	"testing"

	"datasync/internal/state"
)

func TestCursorAggregationKeepsSnapshotDuringConcurrentAdvance(t *testing.T) {
	registry := state.NewCursorRegistry()
	registry.Set("shard-a", 20)
	registry.Set("shard-b", 30)
	ready := make(chan struct{})
	resume := make(chan struct{})
	aggregator := CursorAggregator{Registry: registry, Ready: ready, Continue: resume}
	result := aggregator.SumAsync()
	<-ready
	registry.Set("shard-c", 1000)
	close(resume)
	if got := <-result; got != 50 {
		t.Fatalf("in-flight cursor aggregate changed with a later advance: %d", got)
	}
	nextAggregator := CursorAggregator{Registry: registry}
	next := nextAggregator.SumAsync()
	if got := <-next; got != 1050 {
		t.Fatalf("next cursor aggregate did not include the completed advance: %d", got)
	}
}
