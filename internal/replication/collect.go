package replication

import (
	"context"

	"datasync/internal/stream"
)

// Collect drains the source change stream into a slice. It reads every
// non-failing event from results before reading the terminal error, which
// is the order Produce guarantees. When the stream reports an error the
// events that did succeed are still returned alongside it, so a partial
// batch is never lost.
func Collect(ctx context.Context, items []stream.Event, failAt int) ([]stream.Event, error) {
	results, errorsOut := stream.Produce(ctx, items, failAt)
	collected := make([]stream.Event, 0, len(items))
	for item := range results {
		collected = append(collected, item)
	}
	if err := <-errorsOut; err != nil {
		return collected, err
	}
	return collected, nil
}
