CREATE TABLE IF NOT EXISTS workspace_records (
    kind TEXT NOT NULL CHECK (kind IN ('projects', 'research', 'settings')),
    id TEXT NOT NULL,
    payload JSONB NOT NULL CHECK (jsonb_typeof(payload) = 'object'),
    version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (kind, id)
);

CREATE INDEX IF NOT EXISTS workspace_records_payload_idx ON workspace_records USING GIN(payload);
