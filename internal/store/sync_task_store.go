package store

import (
	"datasync/internal/model"
)

func (s *MemoryStore) CreateSyncTask(t *model.SyncTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.syncTasks {
		if exist.Name == t.Name {
			return ErrConflict
		}
	}
	s.syncTasks[t.ID] = t
	return nil
}

func (s *MemoryStore) GetSyncTask(id string) (*model.SyncTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.syncTasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) GetSyncTaskByName(name string) (*model.SyncTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, t := range s.syncTasks {
		if t.Name == name {
			return t, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListSyncTasks() []*model.SyncTask {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.SyncTask, 0, len(s.syncTasks))
	for _, t := range s.syncTasks {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateSyncTask(t *model.SyncTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.syncTasks[t.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.syncTasks {
		if exist.ID != t.ID && exist.Name == t.Name {
			return ErrConflict
		}
	}
	s.syncTasks[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteSyncTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.syncTasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.syncTasks, id)
	return nil
}
