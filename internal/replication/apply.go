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
		tx.Stage(value)
		err := connector.Normalize(s.Client.Apply(ctx, value))
		if err == nil {
			tx.Commit()
			return nil
		}
		if connector.IsKind(err, connector.ErrorRejected) {
			tx.Rollback()
			return err
		}
		if attempt == attempts-1 {
			tx.Commit()
			return err
		}
	}
	return nil
}
