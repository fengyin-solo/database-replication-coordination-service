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
		for attempt := 0; attempt < 8; attempt++ {
			if err := connector.DetachedCall(d.Client, ctx, stream); err == nil {
				return
			}
			<-time.After(d.delay)
		}
	}()
}

func (d *RetryDispatcher) Wait() { d.wg.Wait() }
