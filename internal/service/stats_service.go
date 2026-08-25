package service

import (
	"sort"
)

type OverviewStats struct {
	TaskTotal          int            `json:"task_total"`
	TaskByStatus       map[string]int `json:"task_by_status"`
	SourceTotal        int            `json:"source_total"`
	TargetTotal        int            `json:"target_total"`
	RecordTotal        int            `json:"record_total"`
	RecordByOperation  map[string]int `json:"record_by_operation"`
	RowsProcessedTotal int64          `json:"rows_processed_total"`
}

func (s *Service) GetOverviewStats() (*OverviewStats, error) {
	tasks := s.store.ListSyncTasks()
	sources := s.store.ListSources()
	targets := s.store.ListTargets()
	records := s.store.ListChangeRecords()
	checkpoints := s.store.ListCheckpoints()

	stats := &OverviewStats{
		TaskTotal:          len(tasks),
		TaskByStatus:       make(map[string]int),
		SourceTotal:        len(sources),
		TargetTotal:        len(targets),
		RecordTotal:        len(records),
		RecordByOperation:  make(map[string]int),
		RowsProcessedTotal: 0,
	}

	for _, t := range tasks {
		stats.TaskByStatus[t.Status]++
	}
	for _, r := range records {
		stats.RecordByOperation[r.Operation]++
	}
	for _, cp := range checkpoints {
		stats.RowsProcessedTotal += cp.RowsProcessed
	}

	return stats, nil
}

type TaskTopItem struct {
	TaskID        string `json:"task_id"`
	TaskName      string `json:"task_name"`
	RowsProcessed int64  `json:"rows_processed"`
}

func (s *Service) GetTasksTop(n int) ([]TaskTopItem, error) {
	if n <= 0 {
		n = 10
	}
	checkpoints := s.store.ListCheckpoints()
	tasks := s.store.ListSyncTasks()
	taskNameMap := make(map[string]string, len(tasks))
	for _, t := range tasks {
		taskNameMap[t.ID] = t.Name
	}

	items := make([]TaskTopItem, 0, len(checkpoints))
	for _, cp := range checkpoints {
		items = append(items, TaskTopItem{
			TaskID:        cp.TaskID,
			TaskName:      taskNameMap[cp.TaskID],
			RowsProcessed: cp.RowsProcessed,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].RowsProcessed > items[j].RowsProcessed
	})

	if len(items) > n {
		items = items[:n]
	}
	return items, nil
}
