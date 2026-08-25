package store

import (
	"datasync/internal/model"
)

func (s *MemoryStore) CreateCheckpoint(c *model.Checkpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.checkpoints {
		if exist.TaskID == c.TaskID {
			return ErrConflict
		}
	}
	s.checkpoints[c.ID] = c
	return nil
}

func (s *MemoryStore) GetCheckpoint(id string) (*model.Checkpoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.checkpoints[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) GetCheckpointByTaskID(taskID string) (*model.Checkpoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.checkpoints {
		if c.TaskID == taskID {
			return c, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListCheckpoints() []*model.Checkpoint {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Checkpoint, 0, len(s.checkpoints))
	for _, c := range s.checkpoints {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) UpdateCheckpoint(c *model.Checkpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.checkpoints[c.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.checkpoints {
		if exist.ID != c.ID && exist.TaskID == c.TaskID {
			return ErrConflict
		}
	}
	s.checkpoints[c.ID] = c
	return nil
}
