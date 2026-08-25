package store

import (
	"datasync/internal/model"
)

func (s *MemoryStore) CreateChangeRecord(c *model.ChangeRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.changeRecords[c.ID] = c
	return nil
}

func (s *MemoryStore) GetChangeRecord(id string) (*model.ChangeRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.changeRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) ListChangeRecords() []*model.ChangeRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ChangeRecord, 0, len(s.changeRecords))
	for _, c := range s.changeRecords {
		list = append(list, c)
	}
	return list
}
