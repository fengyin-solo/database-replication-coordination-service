package replication

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"datasync/internal/stream"
)

func TestCollectReturnsOnSourceFailureAndKeepsNormalEvents(t *testing.T) {
	items := []stream.Event{{ID: "event-a"}, {ID: "event-b"}, {ID: "event-c"}}
	completed := make(chan error, 1)
	go func() {
		_, err := Collect(context.Background(), items, 1)
		completed <- err
	}()
	select {
	case err := <-completed:
		if !errors.Is(err, stream.ErrSourceRead) {
			t.Fatalf("source failure returned the wrong error: %v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("collector hung after the source stream failed")
	}
	got, err := Collect(context.Background(), items, -1)
	if err != nil || !reflect.DeepEqual(got, items) {
		t.Fatalf("normal source stream lost events: got=%v err=%v", got, err)
	}
}
