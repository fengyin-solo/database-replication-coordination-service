// Package stream coordinates finite streams of replication records.
package stream

import (
	"context"
	"errors"
)

var ErrSourceRead = errors.New("source stream read failed")

type Event struct{ ID string }

// Produce streams items to results. A failure at index failAt is recorded
// but does not abort the stream: every non-failing item is still delivered
// so the caller can collect the whole batch. At most one terminal error is
// reported on errorsOut; it is buffered (capacity 1) and sent before the
// channels are closed, so the producer never blocks on the consumer's read
// ordering.
//
// This fixes a deadlock that existed when errorsOut was unbuffered and the
// producer returned early on failure: the producer tried to hand off the
// error while the consumer was still ranging over results, and each side
// was waiting on the other.
func Produce(ctx context.Context, items []Event, failAt int) (<-chan Event, <-chan error) {
	results := make(chan Event)
	// Capacity 1: the single terminal error can be buffered without the
	// producer having to wait for a reader, breaking the old circular wait.
	errorsOut := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(errorsOut)
		var firstErr error
		for index, item := range items {
			if index == failAt && firstErr == nil {
				// Record the failure but keep streaming so the caller can
				// still collect the remaining non-failing events.
				firstErr = ErrSourceRead
				continue
			}
			// Check the context before sending. Relying on select alone is
			// racy: when ctx is already cancelled, select may still pick the
			// send branch, silently dropping the cancellation.
			if err := ctx.Err(); err != nil {
				if firstErr == nil {
					firstErr = err
				}
				break
			}
			select {
			case results <- item:
			case <-ctx.Done():
				if firstErr == nil {
					firstErr = ctx.Err()
				}
				break
			}
		}
		// Single exit point: send the (buffered) terminal error before the
		// deferred closes run. Never blocks — errorsOut has capacity 1 and
		// carries at most one value.
		if firstErr != nil {
			errorsOut <- firstErr
		}
	}()
	return results, errorsOut
}
