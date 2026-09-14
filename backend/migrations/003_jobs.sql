CREATE TABLE jobs (
    id UUID PRIMARY KEY,

    run_id UUID NOT NULL UNIQUE,

    status TEXT NOT NULL,

    attempts INTEGER NOT NULL DEFAULT 0,

    available_at TIMESTAMPTZ NOT NULL,

    claimed_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT jobs_run_fk
        FOREIGN KEY (run_id)
        REFERENCES runs(id)
        ON DELETE CASCADE,

    CONSTRAINT jobs_attempts_non_negative
        CHECK (attempts >= 0),

    CONSTRAINT jobs_status_valid
        CHECK (
            status IN (
                'QUEUED',
                'CLAIMED',
                'COMPLETED',
                'FAILED'
            )
        )
);

CREATE INDEX idx_jobs_status_available_at
    ON jobs(status, available_at);