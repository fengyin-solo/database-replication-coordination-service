package replication

import "datasync/internal/state"

type AttemptService struct{ Store *state.AttemptStore }

func (s *AttemptService) Begin(taskID string) state.Attempt {
	attempt := state.Attempt{TaskID: taskID, Version: 1, Status: "running", Operation: taskID + ":apply"}
	s.Store.ApplyOnce(attempt.Operation)
	s.Store.Update(attempt)
	return attempt
}

func (s *AttemptService) Retry(first state.Attempt) state.Attempt {
	retry := first
	retry.Version++
	retry.Status = "succeeded"
	// Retry 是对同一操作的重新应用，复用首次 Operation：首次尝试
	// 已经在 ApplyOnce 里计入过该操作，因此重试时 ApplyOnce 会返回
	// false，外部变更不会被重复计入。只有首次尝试未真正计入（例如
	// 在 Begin 中被拒绝）时，重试才会补计一次。
	if s.Store.ApplyOnce(retry.Operation) {
		s.Store.Update(retry)
	} else {
		// 操作已被首次尝试计入，仅推进状态，不重复计入。
		if !s.Store.Update(retry) {
			retry.Status = first.Status
			retry.Version = first.Version
		}
	}
	return retry
}

// LateCallback 处理首次尝试的延迟回调：只有当当前 attempt 仍停留在
// 首次尝试的旧版本（尚未被重试推进到成功态）时，才回退为 running。
// 重试已经把版本推进到更高且状态为 succeeded，Update 会因版本更老
// 而拒绝覆盖，函数返回 false，任务保持成功状态。
func (s *AttemptService) LateCallback(first state.Attempt) bool {
	first.Status = "running"
	return s.Store.Update(first)
}
