package replication

import (
	"strings"
	"testing"

	"datasync/internal/state"
)

func TestFailedBatchReleasesLeasesAndPublishesNoSuccess(t *testing.T) {
	pool := state.NewLeasePool(1)
	processor := BatchProcessor{Pool: pool}
	result, err := processor.Process([]string{"change-a", "change-b", "change-c"}, 2)
	if err == nil || !strings.Contains(err.Error(), "item 2") {
		t.Fatalf("batch failure was lost or replaced: %v", err)
	}
	if result.Committed || result.Audited {
		t.Fatalf("failed batch published a success state: %#v", result)
	}
	if pool.InUse() != 0 {
		t.Fatalf("failed batch left replication leases in use: %d", pool.InUse())
	}
	next, err := processor.Process([]string{"change-d"}, -1)
	if err != nil || !next.Committed || !next.Audited {
		t.Fatalf("next batch could not use the released lease: result=%#v err=%v", next, err)
	}
}
