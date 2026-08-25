package replication

import "datasync/internal/state"

type AuditRecord struct {
	Tenant string
	Labels []string
}

type AuditService struct{ Pool *state.RequestPool }

func (s *AuditService) Schedule(tenant string, labels []string, release <-chan struct{}) <-chan AuditRecord {
	result := make(chan AuditRecord, 1)
	request := s.Pool.Acquire(tenant, labels)
	s.Pool.Release(request)
	go func() {
		<-release
		result <- AuditRecord{Tenant: request.Tenant, Labels: request.Labels}
		close(result)
	}()
	return result
}
