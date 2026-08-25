// Package connector defines the remote calls used by the replication engine.
package connector

import (
	"context"
	"errors"
	"fmt"
)

type ErrorKind string

const (
	ErrorTemporary ErrorKind = "temporary"
	ErrorRejected  ErrorKind = "rejected"
)

type RemoteError struct {
	Kind ErrorKind
	Err  error
}

func (e *RemoteError) Error() string { return fmt.Sprintf("remote %s: %v", e.Kind, e.Err) }
func (e *RemoteError) Unwrap() error { return e.Err }

func IsKind(err error, kind ErrorKind) bool {
	var remote *RemoteError
	return errors.As(err, &remote) && remote.Kind == kind
}

// Normalize keeps the typed error available to retry and transaction policy.
func Normalize(err error) error {
	if err == nil {
		return nil
	}
	var remote *RemoteError
	if errors.As(err, &remote) {
		return fmt.Errorf("connector request: %w", err)
	}
	return err
}

type Fetcher interface {
	Fetch(context.Context, string) ([]byte, error)
}

func RequestContext(context.Context) context.Context { return context.Background() }

type Caller interface {
	Call(context.Context, string) error
}

type ApplyClient interface {
	Apply(context.Context, string) error
}
