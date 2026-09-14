package ports

import (
	"context"
	"agentrun/internal/domain"
)

type RunRepository interface {
	Create(ctx context.Context, run *domain.Run) error
	FindByID(ctx context.Context, id string) (*domain.Run, error)
	Update(ctx context.Context, run *domain.Run) error
}