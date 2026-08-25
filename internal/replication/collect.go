package replication

import (
	"context"

	"datasync/internal/stream"
)

func Collect(ctx context.Context, items []stream.Event, failAt int) ([]stream.Event, error) {
	results, errorsOut := stream.Produce(ctx, items, failAt)
	collected := make([]stream.Event, 0, len(items))
	for results != nil || errorsOut != nil {
		select {
		case item, ok := <-results:
			if !ok {
				results = nil
				continue
			}
			collected = append(collected, item)
		case err, ok := <-errorsOut:
			if !ok {
				errorsOut = nil
				continue
			}
			if err != nil {
				return nil, err
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return collected, nil
}
