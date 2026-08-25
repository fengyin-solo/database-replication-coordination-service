// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"datasync/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// SyncTask
	CreateSyncTask(t *model.SyncTask) error
	GetSyncTask(id string) (*model.SyncTask, error)
	GetSyncTaskByName(name string) (*model.SyncTask, error)
	ListSyncTasks() []*model.SyncTask
	UpdateSyncTask(t *model.SyncTask) error
	DeleteSyncTask(id string) error

	// Source
	CreateSource(s *model.Source) error
	GetSource(id string) (*model.Source, error)
	GetSourceByName(name string) (*model.Source, error)
	ListSources() []*model.Source
	UpdateSource(s *model.Source) error
	DeleteSource(id string) error

	// Target
	CreateTarget(t *model.Target) error
	GetTarget(id string) (*model.Target, error)
	GetTargetByName(name string) (*model.Target, error)
	ListTargets() []*model.Target
	UpdateTarget(t *model.Target) error
	DeleteTarget(id string) error

	// ChangeRecord（只读追加）
	CreateChangeRecord(c *model.ChangeRecord) error
	GetChangeRecord(id string) (*model.ChangeRecord, error)
	ListChangeRecords() []*model.ChangeRecord

	// Checkpoint
	CreateCheckpoint(c *model.Checkpoint) error
	GetCheckpoint(id string) (*model.Checkpoint, error)
	GetCheckpointByTaskID(taskID string) (*model.Checkpoint, error)
	ListCheckpoints() []*model.Checkpoint
	UpdateCheckpoint(c *model.Checkpoint) error
}
