package store

import (
	"sync"

	"datasync/internal/model"
)

type MemoryStore struct {
	mu            sync.RWMutex
	syncTasks     map[string]*model.SyncTask
	sources       map[string]*model.Source
	targets       map[string]*model.Target
	changeRecords map[string]*model.ChangeRecord
	checkpoints   map[string]*model.Checkpoint
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		syncTasks:     make(map[string]*model.SyncTask),
		sources:       make(map[string]*model.Source),
		targets:       make(map[string]*model.Target),
		changeRecords: make(map[string]*model.ChangeRecord),
		checkpoints:   make(map[string]*model.Checkpoint),
	}
}

var _ Store = (*MemoryStore)(nil)
