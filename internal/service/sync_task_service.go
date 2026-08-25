package service

import (
	"sort"
	"time"

	"datasync/internal/model"
	"datasync/pkg/idgen"
)

func (s *Service) CreateSyncTask(input model.SyncTask) (*model.SyncTask, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetSource(input.SourceID); err != nil {
		return nil, model.NewValidationError("source_id", "指定的源不存在")
	}
	if _, err := s.store.GetTarget(input.TargetID); err != nil {
		return nil, model.NewValidationError("target_id", "指定的目标不存在")
	}
	now := time.Now().UTC()
	t := &model.SyncTask{
		ID:          idgen.Hex(),
		Name:        input.Name,
		SourceID:    input.SourceID,
		TargetID:    input.TargetID,
		Mode:        input.Mode,
		IntervalSec: input.IntervalSec,
		Status:      model.TaskStatusIdle,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateSyncTask(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) ListSyncTasks(filter model.SyncTaskFilter, page, size int) ([]*model.SyncTask, int, error) {
	all := s.store.ListSyncTasks()
	matched := make([]*model.SyncTask, 0, len(all))
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
		return []*model.SyncTask{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetSyncTask(id string) (*model.SyncTask, error) {
	return s.store.GetSyncTask(id)
}

func (s *Service) UpdateSyncTask(id string, input model.SyncTask) (*model.SyncTask, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	t, err := s.store.GetSyncTask(id)
	if err != nil {
		return nil, err
	}
	if input.SourceID != "" && input.SourceID != t.SourceID {
		if _, err := s.store.GetSource(input.SourceID); err != nil {
			return nil, model.NewValidationError("source_id", "指定的源不存在")
		}
		t.SourceID = input.SourceID
	}
	if input.TargetID != "" && input.TargetID != t.TargetID {
		if _, err := s.store.GetTarget(input.TargetID); err != nil {
			return nil, model.NewValidationError("target_id", "指定的目标不存在")
		}
		t.TargetID = input.TargetID
	}
	t.Name = input.Name
	if input.Mode != "" {
		t.Mode = input.Mode
	}
	if input.IntervalSec > 0 {
		t.IntervalSec = input.IntervalSec
	}
	if input.Status != "" && input.Status != t.Status {
		if !model.CanTransition(t.Status, input.Status) {
			return nil, model.NewValidationError("status", "状态转换不合法")
		}
		t.Status = input.Status
	}
	t.UpdatedAt = time.Now().UTC()
	if err := s.store.UpdateSyncTask(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteSyncTask(id string) error {
	return s.store.DeleteSyncTask(id)
}

func (s *Service) StartSyncTask(id string) (*model.SyncTask, error) {
	t, err := s.store.GetSyncTask(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransition(t.Status, model.TaskStatusRunning) {
		return nil, model.NewValidationError("status", "当前状态无法启动")
	}
	t.Status = model.TaskStatusRunning
	now := time.Now().UTC()
	t.LastRunAt = &now
	t.UpdatedAt = now
	if err := s.store.UpdateSyncTask(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) PauseSyncTask(id string) (*model.SyncTask, error) {
	t, err := s.store.GetSyncTask(id)
	if err != nil {
		return nil, err
	}
	if t.Status == model.TaskStatusRunning {
		if !model.CanTransition(t.Status, model.TaskStatusIdle) {
			return nil, model.NewValidationError("status", "当前状态无法暂停")
		}
		t.Status = model.TaskStatusIdle
	} else if t.Status == model.TaskStatusPaused {
		return t, nil
	} else {
		if !model.CanTransition(t.Status, model.TaskStatusPaused) {
			return nil, model.NewValidationError("status", "当前状态无法暂停")
		}
		t.Status = model.TaskStatusPaused
	}
	t.UpdatedAt = time.Now().UTC()
	if err := s.store.UpdateSyncTask(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) FailSyncTask(id string) (*model.SyncTask, error) {
	t, err := s.store.GetSyncTask(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransition(t.Status, model.TaskStatusFailed) {
		return nil, model.NewValidationError("status", "当前状态无法标记失败")
	}
	t.Status = model.TaskStatusFailed
	t.UpdatedAt = time.Now().UTC()
	if err := s.store.UpdateSyncTask(t); err != nil {
		return nil, err
	}
	return t, nil
}
