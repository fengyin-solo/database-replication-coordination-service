package store

import (
	"datasync/internal/model"
)

func (s *MemoryStore) CreateTarget(t *model.Target) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.targets {
		if exist.Name == t.Name {
			return ErrConflict
		}
	}
	s.targets[t.ID] = t
	return nil
}

func (s *MemoryStore) GetTarget(id string) (*model.Target, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.targets[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) GetTargetByName(name string) (*model.Target, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, t := range s.targets {
		if t.Name == name {
			return t, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListTargets() []*model.Target {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Target, 0, len(s.targets))
	for _, t := range s.targets {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateTarget(t *model.Target) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.targets[t.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.targets {
		if exist.ID != t.ID && exist.Name == t.Name {
			return ErrConflict
		}
	}
	s.targets[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteTarget(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.targets[id]; !ok {
		return ErrNotFound
	}
	delete(s.targets, id)
	return nil
}
