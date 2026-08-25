package model

import (
	"strings"
	"time"
)

const (
	SourceTypeMySQL    = "mysql"
	SourceTypePostgres = "postgres"
	SourceTypeMongo    = "mongo"
	SourceTypeAPI      = "api"

	SourceStatusActive   = "active"
	SourceStatusInactive = "inactive"
)

type Source struct {
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

func (s *Source) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	s.Host = strings.TrimSpace(s.Host)
	s.Database = strings.TrimSpace(s.Database)
	if s.Name == "" {
		return NewValidationError("name", "源名称不能为空")
	}
	if s.Host == "" {
		return NewValidationError("host", "主机地址不能为空")
	}
	if s.Port <= 0 {
		return NewValidationError("port", "端口号必须大于 0")
	}
	if s.Type == "" {
		return NewValidationError("type", "源类型不能为空")
	}
	if s.Type != SourceTypeMySQL && s.Type != SourceTypePostgres && s.Type != SourceTypeMongo && s.Type != SourceTypeAPI {
		return NewValidationError("type", "源类型不合法，可选 mysql/postgres/mongo/api")
	}
	if s.Status == "" {
		s.Status = SourceStatusActive
	}
	if s.Status != SourceStatusActive && s.Status != SourceStatusInactive {
		return NewValidationError("status", "源状态不合法")
	}
	return nil
}

type SourceFilter struct {
	Type    string
	Status  string
	Keyword string
}

func (f SourceFilter) Match(s *Source) bool {
	if f.Type != "" && s.Type != f.Type {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(s.Name), k) && !strings.Contains(strings.ToLower(s.Host), k) {
			return false
		}
	}
	return true
}
