package replication

import (
	"context"
	"sync"
	"time"

	"datasync/internal/connector"
)

type RetryDispatcher struct {
	Client connector.Caller
	delay  time.Duration
	wg     sync.WaitGroup
}

func NewRetryDispatcher(client connector.Caller, delay time.Duration) *RetryDispatcher {
	return &RetryDispatcher{Client: client, delay: delay}
}

func (d *RetryDispatcher) Start(ctx context.Context, stream string) {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		for {
			if err := d.Client.Call(ctx, stream); err == nil {
				return
			}
			timer := time.NewTimer(d.delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
			}
		}
	}()
}

func (d *RetryDispatcher) Wait() { d.wg.Wait() }
