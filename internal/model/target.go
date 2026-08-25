package model

import (
	"strings"
	"time"
)

const (
	TargetTypeMySQL    = "mysql"
	TargetTypePostgres = "postgres"
	TargetTypeMongo    = "mongo"
	TargetTypeAPI      = "api"

	TargetStatusActive   = "active"
	TargetStatusInactive = "inactive"
)

type Target struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Database  string    `json:"database"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *Target) Validate() error {
	t.Name = strings.TrimSpace(t.Name)
	t.Host = strings.TrimSpace(t.Host)
	t.Database = strings.TrimSpace(t.Database)
	if t.Name == "" {
		return NewValidationError("name", "目标名称不能为空")
	}
	if t.Host == "" {
		return NewValidationError("host", "主机地址不能为空")
	}
	if t.Port <= 0 {
		return NewValidationError("port", "端口号必须大于 0")
	}
	if t.Type == "" {
		return NewValidationError("type", "目标类型不能为空")
	}
	if t.Type != TargetTypeMySQL && t.Type != TargetTypePostgres && t.Type != TargetTypeMongo && t.Type != TargetTypeAPI {
		return NewValidationError("type", "目标类型不合法，可选 mysql/postgres/mongo/api")
	}
	if t.Status == "" {
		t.Status = TargetStatusActive
	}
	if t.Status != TargetStatusActive && t.Status != TargetStatusInactive {
		return NewValidationError("status", "目标状态不合法")
	}
	return nil
}

type TargetFilter struct {
	Type    string
	Status  string
	Keyword string
}

func (f TargetFilter) Match(t *Target) bool {
	if f.Type != "" && t.Type != f.Type {
		return false
	}
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.Name), k) && !strings.Contains(strings.ToLower(t.Host), k) {
			return false
		}
	}
	return true
}
