package store

import (
	"datasync/internal/model"
)

func (s *MemoryStore) CreateSource(src *model.Source) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.sources {
		if exist.Name == src.Name {
			return ErrConflict
		}
	}
	s.sources[src.ID] = src
	return nil
}

func (s *MemoryStore) GetSource(id string) (*model.Source, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src, ok := s.sources[id]
	if !ok {
		return nil, ErrNotFound
	}
	return src, nil
}

func (s *MemoryStore) GetSourceByName(name string) (*model.Source, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, src := range s.sources {
		if src.Name == name {
			return src, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListSources() []*model.Source {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Source, 0, len(s.sources))
	for _, src := range s.sources {
		list = append(list, src)
	}
	return list
}

func (s *MemoryStore) UpdateSource(src *model.Source) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sources[src.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.sources {
		if exist.ID != src.ID && exist.Name == src.Name {
			return ErrConflict
		}
	}
	s.sources[src.ID] = src
	return nil
}

func (s *MemoryStore) DeleteSource(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sources[id]; !ok {
		return ErrNotFound
	}
	delete(s.sources, id)
	return nil
}
