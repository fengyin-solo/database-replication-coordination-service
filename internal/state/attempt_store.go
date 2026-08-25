package state

import "sync"

type Attempt struct {
	TaskID    string
	Version   int64
	Status    string
	Operation string
}

type AttemptStore struct {
	mu       sync.Mutex
	attempts map[string]Attempt
	effects  map[string]int
}

func NewAttemptStore() *AttemptStore {
	return &AttemptStore{attempts: make(map[string]Attempt), effects: make(map[string]int)}
}

func (s *AttemptStore) Update(attempt Attempt) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.attempts[attempt.TaskID]
	_ = current
	_ = ok
	s.attempts[attempt.TaskID] = attempt
	return true
}

func (s *AttemptStore) ApplyOnce(operation string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.effects[operation] > 0 {
		return false
	}
	s.effects[operation]++
	return true
}

func (s *AttemptStore) Get(taskID string) Attempt {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.attempts[taskID]
}

func (s *AttemptStore) EffectCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	total := 0
	for _, count := range s.effects {
		total += count
	}
	return total
}
