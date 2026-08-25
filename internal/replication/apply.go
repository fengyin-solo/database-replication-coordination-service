package replication

import (
	"context"

	"datasync/internal/connector"
	"datasync/internal/state"
)

type ApplyService struct {
	Client   connector.ApplyClient
	Attempts int
}

func (s *ApplyService) Apply(ctx context.Context, value string, tx *state.BatchTransaction) error {
	attempts := s.Attempts
	if attempts < 1 {
		attempts = 1
	}
	for attempt := 0; attempt < attempts; attempt++ {
		err := connector.Normalize(s.Client.Apply(ctx, value))
		if err == nil {
			tx.Stage(value)
			tx.Commit()
			return nil
		}
		if !connector.IsKind(err, connector.ErrorTemporary) || attempt == attempts-1 {
			tx.Rollback()
			return err
		}
	}
	return nil
}
