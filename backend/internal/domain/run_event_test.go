package domain

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewRunStartedEvent(t *testing.T) {
	now := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)

	event := NewRunStartedEvent("run-1", now)

	if event.RunID != "run-1" {
		t.Fatalf("expected RunID %q, got %q", "run-1", event.RunID)
	}

	if event.Type != RunEventStarted {
		t.Fatalf("expected type %q, got %q", RunEventStarted, event.Type)
	}

	if !event.CreatedAt.Equal(now) {
		t.Fatalf("expected CreatedAt %v, got %v", now, event.CreatedAt)
	}

	if string(event.Payload) != "{}" {
		t.Fatalf("expected empty payload {}, got %s", string(event.Payload))
	}
}

func TestNewAgentMessageEvent(t *testing.T) {
	now := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)

	event, err := NewAgentMessageEvent(
		"run-1",
		"Inspecting repository",
		now,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if event.Type != RunEventAgentMessage {
		t.Fatalf("expected type %q, got %q", RunEventAgentMessage, event.Type)
	}

	var payload AgentMessagePayload

	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}

	if payload.Message != "Inspecting repository" {
		t.Fatalf(
			"expected message %q, got %q",
			"Inspecting repository",
			payload.Message,
		)
	}
}

func TestNewCommandCompletedEvent(t *testing.T) {
	now := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)

	event, err := NewCommandCompletedEvent(
		"run-1",
		"npm test",
		0,
		4200,
		now,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var payload CommandCompletedPayload

	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}

	if payload.Command != "npm test" {
		t.Fatalf("expected command %q, got %q", "npm test", payload.Command)
	}

	if payload.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", payload.ExitCode)
	}

	if payload.DurationMs != 4200 {
		t.Fatalf("expected duration 4200, got %d", payload.DurationMs)
	}
}