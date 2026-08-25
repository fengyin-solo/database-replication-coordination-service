package replication

import "datasync/internal/state"

type AttemptService struct{ Store *state.AttemptStore }

func (s *AttemptService) Begin(taskID string) state.Attempt {
	attempt := state.Attempt{TaskID: taskID, Version: 1, Status: "running", Operation: taskID + ":apply"}
	s.Store.Update(attempt)
	return attempt
}

func (s *AttemptService) Retry(first state.Attempt) state.Attempt {
	retry := first
	retry.Version++
	retry.Status = "succeeded"
	s.Store.ApplyOnce(retry.Operation)
	s.Store.Update(retry)
	return retry
}

func (s *AttemptService) LateCallback(first state.Attempt) bool {
	first.Status = "running"
	return s.Store.Update(first)
}
