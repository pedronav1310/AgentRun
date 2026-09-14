package domain

import (
	"encoding/json" 
	"time"
)

type RunEventType string

const (
	RunEventCreated   RunEventType = "run.created"
	RunEventStarted   RunEventType = "run.started"
	RunEventCompleted RunEventType = "run.completed"
	RunEventFailed    RunEventType = "run.failed"
	RunEventCancelled RunEventType = "run.cancelled"

	RunEventAgentMessage RunEventType = "agent.message"

	RunEventCommandStarted   RunEventType = "command.started"
	RunEventCommandCompleted RunEventType = "command.completed"
)

type RunEvent struct {
	ID        int64
	RunID     string
	Type      RunEventType
	Payload   json.RawMessage
	CreatedAt time.Time
}

type AgentMessagePayload struct {
	Message string `json:"message"`
}

type CommandStartedPayload struct {
	Command string `json:"command"`
	Cwd     string `json:"cwd,omitempty"`
}

type CommandCompletedPayload struct {
	Command    string `json:"command"`
	ExitCode   int    `json:"exitCode"`
	DurationMs int64  `json:"durationMs"`
}

type RunFailedPayload struct {
	Message string `json:"message"`
}

func NewRunCreatedEvent(runID string, now time.Time) RunEvent {
	return RunEvent{
		RunID:     runID,
		Type:      RunEventCreated,
		Payload:   emptyPayload(),
		CreatedAt: now,
	}
}

func NewRunStartedEvent(runID string, now time.Time) RunEvent {
	return RunEvent{
		RunID:     runID,
		Type:      RunEventStarted,
		Payload:   emptyPayload(),
		CreatedAt: now,
	}
}

func NewRunCompletedEvent(runID string, now time.Time) RunEvent {
	return RunEvent{
		RunID:     runID,
		Type:      RunEventCompleted,
		Payload:   emptyPayload(),
		CreatedAt: now,
	}
}

func NewRunCancelledEvent(runID string, now time.Time) RunEvent {
	return RunEvent{
		RunID:     runID,
		Type:      RunEventCancelled,
		Payload:   emptyPayload(),
		CreatedAt: now,
	}
}

func NewRunFailedEvent(runID string, message string, now time.Time,) (RunEvent, error) {
	return newRunEvent(
		runID,
		RunEventFailed,
		RunFailedPayload{
			Message: message,
		},
		now,
	)
}

func NewAgentMessageEvent(runID string, message string, now time.Time,) (RunEvent, error) {
	return newRunEvent(
		runID,
		RunEventAgentMessage,
		AgentMessagePayload{
			Message: message,
		},
		now,
	)
}

func NewCommandStartedEvent(runID string, command string, cwd string, now time.Time,) (RunEvent, error) {
	return newRunEvent(
		runID,
		RunEventCommandStarted,
		CommandStartedPayload{
			Command: command,
			Cwd:     cwd,
		},
		now,
	)
}

func NewCommandCompletedEvent(runID string, command string, exitCode int, durationMs int64, now time.Time,) (RunEvent, error) {
	return newRunEvent(
		runID,
		RunEventCommandCompleted,
		CommandCompletedPayload{
			Command:    command,
			ExitCode:   exitCode,
			DurationMs: durationMs,
		},
		now,
	)
}

func newRunEvent[T any](runID string, eventType RunEventType, payload T, now time.Time,) (RunEvent, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return RunEvent{}, err
	}

	return RunEvent{
		RunID:     runID,
		Type:      eventType,
		Payload:   data,
		CreatedAt: now,
	}, nil
}

func emptyPayload() json.RawMessage {
	return json.RawMessage(`{}`)
}