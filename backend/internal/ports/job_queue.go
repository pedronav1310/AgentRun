package ports

import (
	"context"
	"agentrun/internal/domain"
)

type JobQueue interface {
	Enqueue(ctx context.Context, job *domain.Job) error
	Claim(ctx context.Context) (*domain.Job, error)
	Update(ctx context.Context, job *domain.Job) error
}