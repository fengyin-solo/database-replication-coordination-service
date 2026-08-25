package replication

import (
	"context"

	"datasync/internal/stream"
)

func Collect(ctx context.Context, items []stream.Event, failAt int) ([]stream.Event, error) {
	results, errorsOut := stream.Produce(ctx, items, failAt)
	collected := make([]stream.Event, 0, len(items))
	for item := range results {
		collected = append(collected, item)
	}
	if err := <-errorsOut; err != nil {
		return nil, err
	}
	return collected, nil
}
