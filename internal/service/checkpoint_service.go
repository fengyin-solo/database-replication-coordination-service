package service

import (
	"sort"
	"time"

	"datasync/internal/model"
	"datasync/pkg/idgen"
)

func (s *Service) CreateCheckpoint(input model.Checkpoint) (*model.Checkpoint, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetSyncTask(input.TaskID); err != nil {
		return nil, model.NewValidationError("task_id", "指定的任务不存在")
	}
	now := time.Now().UTC()
	c := &model.Checkpoint{
		ID:            idgen.Hex(),
		TaskID:        input.TaskID,
		Position:      input.Position,
		RowsProcessed: input.RowsProcessed,
		UpdatedAt:     now,
	}
	if err := s.store.CreateCheckpoint(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) ListCheckpoints(filter model.CheckpointFilter, page, size int) ([]*model.Checkpoint, int, error) {
	all := s.store.ListCheckpoints()
	matched := make([]*model.Checkpoint, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].UpdatedAt.After(matched[j].UpdatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Checkpoint{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetCheckpointByTaskID(taskID string) (*model.Checkpoint, error) {
	return s.store.GetCheckpointByTaskID(taskID)
}

func (s *Service) AdvanceCheckpoint(taskID string, position string, rowsDelta int64) (*model.Checkpoint, error) {
	c, err := s.store.GetCheckpointByTaskID(taskID)
	if err != nil {
		return nil, err
	}
	if position != "" {
		c.Position = position
	}
	if rowsDelta > 0 {
		c.RowsProcessed += rowsDelta
	}
	c.UpdatedAt = time.Now().UTC()
	if err := s.store.UpdateCheckpoint(c); err != nil {
		return nil, err
	}
	return c, nil
}
