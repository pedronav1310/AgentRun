CREATE TABLE run_events (
    id BIGSERIAL PRIMARY KEY,

    run_id UUID NOT NULL,

    event_type TEXT NOT NULL,

    payload JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT run_events_run_fk
        FOREIGN KEY (run_id)
        REFERENCES runs(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_run_events_run_id_created_at
    ON run_events(run_id, created_at);