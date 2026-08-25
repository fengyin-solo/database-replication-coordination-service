package replication

import (
	"context"

	"datasync/internal/connector"
)

type CatchupService struct{ Fetcher connector.Fetcher }

func (s *CatchupService) Catchup(ctx context.Context, stream string) ([]byte, error) {
	return s.Fetcher.Fetch(connector.RequestContext(ctx), stream)
}
