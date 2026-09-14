package domain

import (
	"errors"
	"time"
)

type RunStatus string

const (
	RunStatusQueued    RunStatus = "QUEUED"
	RunStatusRunning   RunStatus = "RUNNING"
	RunStatusCompleted RunStatus = "COMPLETED"
	RunStatusFailed    RunStatus = "FAILED"
	RunStatusCancelled RunStatus = "CANCELLED"
)

var (
	ErrRunAlreadyStarted      = errors.New("run has already started")
	ErrRunNotRunning          = errors.New("run is not running")
	ErrRunAlreadyTerminal     = errors.New("run is already in a terminal state")
	ErrRunCannotBeCancelled   = errors.New("run cannot be cancelled")
)

type Run struct {
	ID string

	RepositoryOwner string
	RepositoryName  string
	IssueNumber     int

	Agent string

	Status RunStatus

	StartedAt  *time.Time
	FinishedAt *time.Time

	Error *string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewRun(id string, repositoryOwner string, repositoryName string, issueNumber int, agent string, now time.Time,) *Run {
	return &Run{
		ID:              id,
		RepositoryOwner: repositoryOwner,
		RepositoryName:  repositoryName,
		IssueNumber:     issueNumber,
		Agent:           agent,
		Status:          RunStatusQueued,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (r *Run) isTerminal() bool{
	switch r.Status{
	case RunStatusCompleted, RunStatusFailed, RunStatusCancelled: 
		return true
	default:
		return false
	}
}

func (r *Run) Start(now time.Time) error {
	if r.isTerminal() {
		return ErrRunAlreadyTerminal
	}

	if r.Status != RunStatusQueued {
		return ErrRunAlreadyStarted
	}

	r.Status = RunStatusRunning
	r.StartedAt = &now
	r.UpdatedAt = now

	return nil
}

func (r *Run) Complete(now time.Time) error{
	if r.isTerminal(){
		return ErrRunAlreadyTerminal
	}

	if r.Status != RunStatusRunning{
		return ErrRunNotRunning
	}

	r.Status = RunStatusCompleted
	r.FinishedAt = &now
	r.UpdatedAt = now

	return nil
}

func (r *Run) Fail(message string, now time.Time) error {
	if r.isTerminal(){
		return ErrRunAlreadyTerminal
	}

	if r.Status != RunStatusRunning{
		return ErrRunNotRunning
	}

	r.Status = RunStatusFailed
	r.Error = &message
	r.FinishedAt = &now
	r.UpdatedAt = now

	return nil
}

func (r *Run) Cancel(now time.Time) error{
	if r.isTerminal(){
		return ErrRunAlreadyTerminal
	}

	switch r.Status{
	case RunStatusQueued, RunStatusRunning:
		r.Status = RunStatusCancelled
		r.FinishedAt = &now
		r.UpdatedAt = now
		return nil

	default: return ErrRunCannotBeCancelled
	}


}

