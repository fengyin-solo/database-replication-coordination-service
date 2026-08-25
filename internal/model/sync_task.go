package model

import (
	"strings"
	"time"
)

const (
	SyncModeFull        = "full"
	SyncModeIncremental = "incremental"

	TaskStatusIdle    = "idle"
	TaskStatusRunning = "running"
	TaskStatusPaused  = "paused"
	TaskStatusFailed  = "failed"
)

var taskTransitions = map[string]map[string]bool{
	TaskStatusIdle:    {TaskStatusRunning: true, TaskStatusPaused: true},
	TaskStatusRunning: {TaskStatusIdle: true, TaskStatusPaused: true, TaskStatusFailed: true},
	TaskStatusPaused:  {TaskStatusIdle: true},
	TaskStatusFailed:  {TaskStatusIdle: true},
}

func CanTransition(from, to string) bool {
	if m, ok := taskTransitions[from]; ok {
		return m[to]
	}
	return false
}

type SyncTask struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	SourceID    string     `json:"source_id"`
	TargetID    string     `json:"target_id"`
	Mode        string     `json:"mode"`
	IntervalSec int        `json:"interval_sec"`
	Status      string     `json:"status"`
	LastRunAt   *time.Time `json:"last_run_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (t *SyncTask) Validate() error {
	t.Name = strings.TrimSpace(t.Name)
	if t.Name == "" {
		return NewValidationError("name", "任务名称不能为空")
	}
	if t.SourceID == "" {
		return NewValidationError("source_id", "源 ID 不能为空")
	}
	if t.TargetID == "" {
		return NewValidationError("target_id", "目标 ID 不能为空")
	}
	if t.Mode == "" {
		t.Mode = SyncModeFull
	}
	if t.Mode != SyncModeFull && t.Mode != SyncModeIncremental {
		return NewValidationError("mode", "同步模式不合法，可选 full/incremental")
	}
	if t.IntervalSec <= 0 {
		t.IntervalSec = 60
	}
	if t.Status == "" {
		t.Status = TaskStatusIdle
	}
	if t.Status != TaskStatusIdle && t.Status != TaskStatusRunning && t.Status != TaskStatusPaused && t.Status != TaskStatusFailed {
		return NewValidationError("status", "任务状态不合法")
	}
	return nil
}

type SyncTaskFilter struct {
	Mode    string
	Status  string
	Keyword string
}

func (f SyncTaskFilter) Match(t *SyncTask) bool {
	if f.Mode != "" && t.Mode != f.Mode {
		return false
	}
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.Name), k) {
			return false
		}
	}
	return true
}
