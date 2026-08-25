package service

import (
	"sort"
	"time"

	"datasync/internal/model"
	"datasync/pkg/idgen"
)

func (s *Service) CreateChangeRecord(input model.ChangeRecord) (*model.ChangeRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetSyncTask(input.TaskID); err != nil {
		return nil, model.NewValidationError("task_id", "指定的任务不存在")
	}
	now := time.Now().UTC()
	c := &model.ChangeRecord{
		ID:        idgen.Hex(),
		TaskID:    input.TaskID,
		TableName: input.TableName,
		Operation: input.Operation,
		Payload:   input.Payload,
		BinlogPos: input.BinlogPos,
		CreatedAt: now,
	}
	if err := s.store.CreateChangeRecord(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) ListChangeRecords(filter model.ChangeRecordFilter, page, size int) ([]*model.ChangeRecord, int, error) {
	all := s.store.ListChangeRecords()
	matched := make([]*model.ChangeRecord, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ChangeRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetChangeRecord(id string) (*model.ChangeRecord, error) {
	return s.store.GetChangeRecord(id)
}
