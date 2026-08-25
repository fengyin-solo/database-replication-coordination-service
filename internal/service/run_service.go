package service

import (
	"time"

	"datasync/internal/model"
	"datasync/pkg/idgen"
)

func (s *Service) RunSyncTask(taskID string) (*model.SyncTask, error) {
	if taskID == "" {
		return nil, model.NewValidationError("task_id", "任务 ID 不能为空")
	}
	t, err := s.store.GetSyncTask(taskID)
	if err != nil {
		return nil, err
	}
	if t.Status == model.TaskStatusRunning {
		return nil, model.NewValidationError("status", "任务已在运行中")
	}
	if _, err := s.store.GetSource(t.SourceID); err != nil {
		return nil, model.NewValidationError("source_id", "任务关联的源不存在")
	}
	if _, err := s.store.GetTarget(t.TargetID); err != nil {
		return nil, model.NewValidationError("target_id", "任务关联的目标不存在")
	}
	if !model.CanTransition(t.Status, model.TaskStatusRunning) {
		return nil, model.NewValidationError("status", "当前状态无法启动运行")
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

func (s *Service) RecordChangeAndAdvance(taskID, tableName, operation, payload, binlogPos string) (*model.ChangeRecord, error) {
	if taskID == "" {
		return nil, model.NewValidationError("task_id", "任务 ID 不能为空")
	}
	if tableName == "" {
		return nil, model.NewValidationError("table_name", "表名不能为空")
	}
	if operation == "" {
		return nil, model.NewValidationError("operation", "操作类型不能为空")
	}
	if payload == "" {
		return nil, model.NewValidationError("payload", "Payload 不能为空")
	}
	if operation != model.OpInsert && operation != model.OpUpdate && operation != model.OpDelete {
		return nil, model.NewValidationError("operation", "操作类型不合法")
	}
	if _, err := s.store.GetSyncTask(taskID); err != nil {
		return nil, model.NewValidationError("task_id", "指定的任务不存在")
	}
	rec := &model.ChangeRecord{
		ID:        idgen.Hex(),
		TaskID:    taskID,
		TableName: tableName,
		Operation: operation,
		Payload:   payload,
		BinlogPos: binlogPos,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.store.CreateChangeRecord(rec); err != nil {
		return nil, err
	}
	if cp, err := s.store.GetCheckpointByTaskID(taskID); err == nil {
		if binlogPos != "" {
			cp.Position = binlogPos
		}
		cp.RowsProcessed += 1
		cp.UpdatedAt = time.Now().UTC()
		_ = s.store.UpdateCheckpoint(cp)
	}
	return rec, nil
}
