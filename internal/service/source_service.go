package service

import (
	"sort"
	"time"

	"datasync/internal/model"
	"datasync/pkg/idgen"
)

func (s *Service) CreateSource(input model.Source) (*model.Source, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	src := &model.Source{
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
	if err := s.store.CreateSource(src); err != nil {
		return nil, err
	}
	return src, nil
}

func (s *Service) ListSources(filter model.SourceFilter, page, size int) ([]*model.Source, int, error) {
	all := s.store.ListSources()
	matched := make([]*model.Source, 0, len(all))
	for _, src := range all {
		if filter.Match(src) {
			matched = append(matched, src)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Source{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetSource(id string) (*model.Source, error) {
	return s.store.GetSource(id)
}

func (s *Service) UpdateSource(id string, input model.Source) (*model.Source, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	src, err := s.store.GetSource(id)
	if err != nil {
		return nil, err
	}
	src.Name = input.Name
	src.Type = input.Type
	src.Host = input.Host
	src.Port = input.Port
	src.Database = input.Database
	if input.Status != "" {
		src.Status = input.Status
	}
	src.UpdatedAt = time.Now().UTC()
	if err := s.store.UpdateSource(src); err != nil {
		return nil, err
	}
	return src, nil
}

func (s *Service) DeleteSource(id string) error {
	return s.store.DeleteSource(id)
}
