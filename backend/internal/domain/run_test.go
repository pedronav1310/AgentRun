package domain

import (
	"testing"
	"time"
)

func TestNewRunStartsQueued(t *testing.T) {
	now := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)

	run := NewRun(
		"run-1",
		"pedronav1310",
		"prproof",
		28,
		"claude",
		now,
	)

	if run.Status != RunStatusQueued {
		t.Fatalf("expected status %s, got %s", RunStatusQueued, run.Status)
	}

	if !run.CreatedAt.Equal(now) {
		t.Fatalf("expected CreatedAt %v, got %v", now, run.CreatedAt)
	}

	if !run.UpdatedAt.Equal(now) {
		t.Fatalf("expected UpdatedAt %v, got %v", now, run.UpdatedAt)
	}

	if run.StartedAt != nil {
		t.Fatal("expected StartedAt to be nil")
	}

	if run.FinishedAt != nil {
		t.Fatal("expected FinishedAt to be nil")
	}
}

func TestRunStartTransitionsQueuedToRunning(t *testing.T) {
	createdAt := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)
	startedAt := createdAt.Add(time.Minute)

	run := NewRun(
		"run-1",
		"pedronav1310",
		"prproof",
		28,
		"claude",
		createdAt,
	)

	err := run.Start(startedAt)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if run.Status != RunStatusRunning {
		t.Fatalf("expected status %s, got %s", RunStatusRunning, run.Status)
	}

	if run.StartedAt == nil {
		t.Fatal("expected StartedAt to be set")
	}

	if !run.StartedAt.Equal(startedAt) {
		t.Fatalf("expected StartedAt %v, got %v", startedAt, *run.StartedAt)
	}

	if !run.UpdatedAt.Equal(startedAt) {
		t.Fatalf("expected UpdatedAt %v, got %v", startedAt, run.UpdatedAt)
	}
}

func TestRunCompleteTransitionsRunningToCompleted(t *testing.T) {
	createdAt := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)
	startedAt := createdAt.Add(time.Minute)
	finishedAt := startedAt.Add(5 * time.Minute)

	run := NewRun(
		"run-1",
		"pedronav1310",
		"prproof",
		28,
		"claude",
		createdAt,
	)

	if err := run.Start(startedAt); err != nil {
		t.Fatalf("failed to start run: %v", err)
	}

	err := run.Complete(finishedAt)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if run.Status != RunStatusCompleted {
		t.Fatalf("expected status %s, got %s", RunStatusCompleted, run.Status)
	}

	if run.FinishedAt == nil {
		t.Fatal("expected FinishedAt to be set")
	}

	if !run.FinishedAt.Equal(finishedAt) {
		t.Fatalf("expected FinishedAt %v, got %v", finishedAt, *run.FinishedAt)
	}
}

func TestRunFailTransitionsRunningToFailed(t *testing.T) {
	createdAt := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)
	startedAt := createdAt.Add(time.Minute)
	failedAt := startedAt.Add(2 * time.Minute)

	run := NewRun(
		"run-1",
		"pedronav1310",
		"prproof",
		28,
		"claude",
		createdAt,
	)

	if err := run.Start(startedAt); err != nil {
		t.Fatalf("failed to start run: %v", err)
	}

	err := run.Fail("sandbox crashed", failedAt)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if run.Status != RunStatusFailed {
		t.Fatalf("expected status %s, got %s", RunStatusFailed, run.Status)
	}

	if run.Error == nil {
		t.Fatal("expected Error to be set")
	}

	if *run.Error != "sandbox crashed" {
		t.Fatalf("expected error message %q, got %q", "sandbox crashed", *run.Error)
	}

	if run.FinishedAt == nil {
		t.Fatal("expected FinishedAt to be set")
	}
}

func TestRunCanBeCancelledWhileQueued(t *testing.T) {
	now := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)
	cancelledAt := now.Add(time.Minute)

	run := NewRun(
		"run-1",
		"pedronav1310",
		"prproof",
		28,
		"claude",
		now,
	)

	err := run.Cancel(cancelledAt)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if run.Status != RunStatusCancelled {
		t.Fatalf("expected status %s, got %s", RunStatusCancelled, run.Status)
	}
}

func TestRunCanBeCancelledWhileRunning(t *testing.T) {
	now := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)

	run := NewRun(
		"run-1",
		"pedronav1310",
		"prproof",
		28,
		"claude",
		now,
	)

	if err := run.Start(now.Add(time.Minute)); err != nil {
		t.Fatalf("failed to start run: %v", err)
	}

	err := run.Cancel(now.Add(2 * time.Minute))

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if run.Status != RunStatusCancelled {
		t.Fatalf("expected status %s, got %s", RunStatusCancelled, run.Status)
	}
}

func TestRunRejectsInvalidTransitions(t *testing.T) {
	now := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		setup func(*Run)
		act   func(*Run) error
	}{
		{
			name: "queued cannot complete",
			setup: func(run *Run) {},
			act: func(run *Run) error {
				return run.Complete(now.Add(time.Minute))
			},
		},
		{
			name: "queued cannot fail",
			setup: func(run *Run) {},
			act: func(run *Run) error {
				return run.Fail("failure", now.Add(time.Minute))
			},
		},
		{
			name: "completed cannot start again",
			setup: func(run *Run) {
				if err := run.Start(now.Add(time.Minute)); err != nil {
					t.Fatalf("failed to start run: %v", err)
				}

				if err := run.Complete(now.Add(2 * time.Minute)); err != nil {
					t.Fatalf("failed to complete run: %v", err)
				}
			},
			act: func(run *Run) error {
				return run.Start(now.Add(3 * time.Minute))
			},
		},
		{
			name: "failed cannot complete",
			setup: func(run *Run) {
				if err := run.Start(now.Add(time.Minute)); err != nil {
					t.Fatalf("failed to start run: %v", err)
				}

				if err := run.Fail("failure", now.Add(2*time.Minute)); err != nil {
					t.Fatalf("failed to fail run: %v", err)
				}
			},
			act: func(run *Run) error {
				return run.Complete(now.Add(3 * time.Minute))
			},
		},
		{
			name: "cancelled cannot start",
			setup: func(run *Run) {
				if err := run.Cancel(now.Add(time.Minute)); err != nil {
					t.Fatalf("failed to cancel run: %v", err)
				}
			},
			act: func(run *Run) error {
				return run.Start(now.Add(2 * time.Minute))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			run := NewRun(
				"run-1",
				"pedronav1310",
				"prproof",
				28,
				"claude",
				now,
			)

			tt.setup(run)

			err := tt.act(run)

			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}