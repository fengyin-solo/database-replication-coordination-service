package model

import (
	"strings"
	"time"
)

type Checkpoint struct {
	ID            string    `json:"id"`
	TaskID        string    `json:"task_id"`
	Position      string    `json:"position"`
	RowsProcessed int64     `json:"rows_processed"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (c *Checkpoint) Validate() error {
	c.TaskID = strings.TrimSpace(c.TaskID)
	c.Position = strings.TrimSpace(c.Position)
	if c.TaskID == "" {
		return NewValidationError("task_id", "任务 ID 不能为空")
	}
	if c.RowsProcessed < 0 {
		return NewValidationError("rows_processed", "已处理行数不能为负数")
	}
	return nil
}

type CheckpointFilter struct {
	TaskID string
}

func (f CheckpointFilter) Match(c *Checkpoint) bool {
	if f.TaskID != "" && c.TaskID != f.TaskID {
		return false
	}
	return true
}
