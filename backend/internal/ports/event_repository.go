package ports

import (
	"context"
	"agentrun/internal/domain"
)

type EventRepository interface {
	Append(ctx context.Context, event domain.RunEvent) error
	ListByRunID(ctx context.Context, runID string) ([]domain.RunEvent, error)
}