package service

import (
	"sort"
	"time"

	"datasync/internal/model"
	"datasync/pkg/idgen"
)

func (s *Service) CreateTarget(input model.Target) (*model.Target, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	t := &model.Target{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Type:      input.Type,
		Host:      input.Host,
		Port:      input.Port,
		Database:  input.Database,
		Status:    input.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateTarget(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) ListTargets(filter model.TargetFilter, page, size int) ([]*model.Target, int, error) {
	all := s.store.ListTargets()
	matched := make([]*model.Target, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Target{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetTarget(id string) (*model.Target, error) {
	return s.store.GetTarget(id)
}

func (s *Service) UpdateTarget(id string, input model.Target) (*model.Target, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	t, err := s.store.GetTarget(id)
	if err != nil {
		return nil, err
	}
	t.Name = input.Name
	t.Type = input.Type
	t.Host = input.Host
	t.Port = input.Port
	t.Database = input.Database
	if input.Status != "" {
		t.Status = input.Status
	}
	t.UpdatedAt = time.Now().UTC()
	if err := s.store.UpdateTarget(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteTarget(id string) error {
	return s.store.DeleteTarget(id)
}
