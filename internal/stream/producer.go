// Package stream coordinates finite streams of replication records.
package stream

import (
	"context"
	"errors"
)

var ErrSourceRead = errors.New("source stream read failed")

type Event struct{ ID string }

func Produce(ctx context.Context, items []Event, failAt int) (<-chan Event, <-chan error) {
	results := make(chan Event)
	errorsOut := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errorsOut)
		for index, item := range items {
			if index == failAt {
				errorsOut <- ErrSourceRead
				return
			}
			select {
			case results <- item:
			case <-ctx.Done():
				errorsOut <- ctx.Err()
				return
			}
		}
	}()
	return results, errorsOut
}
