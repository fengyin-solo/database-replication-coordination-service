package replication

import (
	"testing"

	"datasync/internal/state"
)

func TestSuccessfulRetryRejectsLateCallbackAndDuplicateApply(t *testing.T) {
	store := state.NewAttemptStore()
	service := AttemptService{Store: store}
	first := service.Begin("replication-a")
	retry := service.Retry(first)
	if retry.Operation != first.Operation {
		t.Fatalf("retry changed the stable apply identity: first=%q retry=%q", first.Operation, retry.Operation)
	}
	if service.LateCallback(first) {
		t.Fatal("late first-attempt callback was accepted after retry success")
	}
	final := store.Get("replication-a")
	if final.Status != "succeeded" || final.Version != 2 {
		t.Fatalf("successful retry regressed to an older state: %#v", final)
	}
	if count := store.EffectCount(); count != 1 {
		t.Fatalf("external apply ran more than once across retry: %d", count)
	}
}
