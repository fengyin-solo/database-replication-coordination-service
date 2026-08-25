package model

import (
	"strings"
	"time"
)

const (
	OpInsert = "insert"
	OpUpdate = "update"
	OpDelete = "delete"
)

type ChangeRecord struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	TableName string    `json:"table_name"`
	Operation string    `json:"operation"`
	Payload   string    `json:"payload"`
	BinlogPos string    `json:"binlog_pos"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *ChangeRecord) Validate() error {
	c.TableName = strings.TrimSpace(c.TableName)
	c.Payload = strings.TrimSpace(c.Payload)
	c.BinlogPos = strings.TrimSpace(c.BinlogPos)
	if c.TaskID == "" {
		return NewValidationError("task_id", "任务 ID 不能为空")
	}
	if c.TableName == "" {
		return NewValidationError("table_name", "表名不能为空")
	}
	if c.Operation == "" {
		return NewValidationError("operation", "操作类型不能为空")
	}
	if c.Operation != OpInsert && c.Operation != OpUpdate && c.Operation != OpDelete {
		return NewValidationError("operation", "操作类型不合法，可选 insert/update/delete")
	}
	if c.Payload == "" {
		return NewValidationError("payload", "Payload 不能为空")
	}
	return nil
}

type ChangeRecordFilter struct {
	TaskID    string
	TableName string
	Operation string
}

func (f ChangeRecordFilter) Match(c *ChangeRecord) bool {
	if f.TaskID != "" && c.TaskID != f.TaskID {
		return false
	}
	if f.TableName != "" && c.TableName != f.TableName {
		return false
	}
	if f.Operation != "" && c.Operation != f.Operation {
		return false
	}
	return true
}
