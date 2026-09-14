package domain

import (
	"testing"
	"time"
)

func TestNewJobStartsQueued(t *testing.T) {
	now := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)

	job := NewJob(
		"job-1",
		"run-1",
		now,
	)

	if job.Status != JobStatusQueued {
		t.Fatalf("expected status %s, got %s", JobStatusQueued, job.Status)
	}

	if job.Attempts != 0 {
		t.Fatalf("expected attempts 0, got %d", job.Attempts)
	}

	if !job.AvailableAt.Equal(now) {
		t.Fatalf("expected AvailableAt %v, got %v", now, job.AvailableAt)
	}

	if !job.CreatedAt.Equal(now) {
		t.Fatalf("expected CreatedAt %v, got %v", now, job.CreatedAt)
	}

	if !job.UpdatedAt.Equal(now) {
		t.Fatalf("expected UpdatedAt %v, got %v", now, job.UpdatedAt)
	}

	if job.ClaimedAt != nil {
		t.Fatal("expected ClaimedAt to be nil")
	}

	if job.CompletedAt != nil {
		t.Fatal("expected CompletedAt to be nil")
	}
}

func TestJobClaimTransitionsQueuedToClaimed(t *testing.T) {
	createdAt := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)
	claimedAt := createdAt.Add(time.Minute)

	job := NewJob(
		"job-1",
		"run-1",
		createdAt,
	)

	err := job.Claim(claimedAt)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if job.Status != JobStatusClaimed {
		t.Fatalf("expected status %s, got %s", JobStatusClaimed, job.Status)
	}

	if job.Attempts != 1 {
		t.Fatalf("expected attempts 1, got %d", job.Attempts)
	}

	if job.ClaimedAt == nil {
		t.Fatal("expected ClaimedAt to be set")
	}

	if !job.ClaimedAt.Equal(claimedAt) {
		t.Fatalf("expected ClaimedAt %v, got %v", claimedAt, *job.ClaimedAt)
	}

	if !job.UpdatedAt.Equal(claimedAt) {
		t.Fatalf("expected UpdatedAt %v, got %v", claimedAt, job.UpdatedAt)
	}
}

func TestJobCompleteTransitionsClaimedToCompleted(t *testing.T) {
	createdAt := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)
	claimedAt := createdAt.Add(time.Minute)
	completedAt := claimedAt.Add(5 * time.Minute)

	job := NewJob(
		"job-1",
		"run-1",
		createdAt,
	)

	if err := job.Claim(claimedAt); err != nil {
		t.Fatalf("failed to claim job: %v", err)
	}

	err := job.Complete(completedAt)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if job.Status != JobStatusCompleted {
		t.Fatalf("expected status %s, got %s", JobStatusCompleted, job.Status)
	}

	if job.CompletedAt == nil {
		t.Fatal("expected CompletedAt to be set")
	}

	if !job.CompletedAt.Equal(completedAt) {
		t.Fatalf(
			"expected CompletedAt %v, got %v",
			completedAt,
			*job.CompletedAt,
		)
	}
}

func TestJobFailTransitionsClaimedToFailed(t *testing.T) {
	createdAt := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)
	claimedAt := createdAt.Add(time.Minute)
	failedAt := claimedAt.Add(2 * time.Minute)

	job := NewJob(
		"job-1",
		"run-1",
		createdAt,
	)

	if err := job.Claim(claimedAt); err != nil {
		t.Fatalf("failed to claim job: %v", err)
	}

	err := job.Fail(failedAt)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if job.Status != JobStatusFailed {
		t.Fatalf("expected status %s, got %s", JobStatusFailed, job.Status)
	}

	if job.CompletedAt == nil {
		t.Fatal("expected CompletedAt to be set")
	}

	if !job.CompletedAt.Equal(failedAt) {
		t.Fatalf(
			"expected CompletedAt %v, got %v",
			failedAt,
			*job.CompletedAt,
		)
	}
}

func TestJobRejectsInvalidTransitions(t *testing.T) {
	now := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		setup func(*Job)
		act   func(*Job) error
	}{
		{
			name:  "queued cannot complete",
			setup: func(job *Job) {},
			act: func(job *Job) error {
				return job.Complete(now.Add(time.Minute))
			},
		},
		{
			name:  "queued cannot fail",
			setup: func(job *Job) {},
			act: func(job *Job) error {
				return job.Fail(now.Add(time.Minute))
			},
		},
		{
			name: "claimed cannot be claimed again",
			setup: func(job *Job) {
				if err := job.Claim(now.Add(time.Minute)); err != nil {
					t.Fatalf("failed to claim job: %v", err)
				}
			},
			act: func(job *Job) error {
				return job.Claim(now.Add(2 * time.Minute))
			},
		},
		{
			name: "completed cannot be claimed",
			setup: func(job *Job) {
				if err := job.Claim(now.Add(time.Minute)); err != nil {
					t.Fatalf("failed to claim job: %v", err)
				}

				if err := job.Complete(now.Add(2 * time.Minute)); err != nil {
					t.Fatalf("failed to complete job: %v", err)
				}
			},
			act: func(job *Job) error {
				return job.Claim(now.Add(3 * time.Minute))
			},
		},
		{
			name: "failed cannot complete",
			setup: func(job *Job) {
				if err := job.Claim(now.Add(time.Minute)); err != nil {
					t.Fatalf("failed to claim job: %v", err)
				}

				if err := job.Fail(now.Add(2 * time.Minute)); err != nil {
					t.Fatalf("failed to fail job: %v", err)
				}
			},
			act: func(job *Job) error {
				return job.Complete(now.Add(3 * time.Minute))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := NewJob(
				"job-1",
				"run-1",
				now,
			)

			tt.setup(job)

			err := tt.act(job)

			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}
