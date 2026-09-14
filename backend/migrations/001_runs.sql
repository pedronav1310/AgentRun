CREATE TABLE runs (
    id UUID PRIMARY KEY,

    repository_owner TEXT NOT NULL,
    repository_name TEXT NOT NULL,
    issue_number INTEGER NOT NULL,

    agent TEXT NOT NULL,

    status TEXT NOT NULL,

    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,

    error TEXT,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT runs_issue_number_positive
        CHECK (issue_number > 0),

    CONSTRAINT runs_status_valid
        CHECK (
            status IN (
                'QUEUED',
                'RUNNING',
                'COMPLETED',
                'FAILED',
                'CANCELLED'
            )
        )
);