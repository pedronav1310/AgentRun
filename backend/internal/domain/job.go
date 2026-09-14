package domain

import (
	"errors"
	"time"
)

type JobStatus string

const (
	JobStatusQueued    JobStatus = "QUEUED"
	JobStatusClaimed   JobStatus = "CLAIMED"
	JobStatusCompleted JobStatus = "COMPLETED"
	JobStatusFailed    JobStatus = "FAILED"
)

var (
	ErrJobNotQueued       = errors.New("job is not queued")
	ErrJobNotClaimed      = errors.New("job is not claimed")
	ErrJobAlreadyTerminal = errors.New("job is already in a terminal state")
)

type Job struct {
	ID    string
	RunID string

	Status   JobStatus
	Attempts int

	AvailableAt time.Time
	ClaimedAt   *time.Time
	CompletedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (j *Job) isTerminal() bool {
	return j.Status == JobStatusCompleted ||
		j.Status == JobStatusFailed
}

func NewJob(id string, runID string, now time.Time) *Job {
	return &Job{
		ID:          id,
		RunID:       runID,
		Status:      JobStatusQueued,
		Attempts:    0,
		AvailableAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (j *Job) Claim(now time.Time) error {
	if j.isTerminal() {
		return ErrJobAlreadyTerminal
	}

	if j.Status != JobStatusQueued {
		return ErrJobNotQueued
	}

	j.Status = JobStatusClaimed
	j.Attempts++
	j.ClaimedAt = &now
	j.UpdatedAt = now

	return nil
}

func (j *Job) Complete(now time.Time) error {
	if j.isTerminal() {
		return ErrJobAlreadyTerminal
	}

	if j.Status != JobStatusClaimed {
		return ErrJobNotClaimed
	}

	j.Status = JobStatusCompleted
	j.CompletedAt = &now
	j.UpdatedAt = now

	return nil
}

func (j *Job) Fail(now time.Time) error {
	if j.isTerminal() {
		return ErrJobAlreadyTerminal
	}

	if j.Status != JobStatusClaimed {
		return ErrJobNotClaimed
	}

	j.Status = JobStatusFailed
	j.CompletedAt = &now
	j.UpdatedAt = now

	return nil
}
