CREATE TABLE entities (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    kind TEXT NOT NULL,
    attributes JSONB NOT NULL DEFAULT '{}'::jsonb,
    version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE sources (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT,
    publisher TEXT,
    url TEXT,
    published_time JSONB,
    retrieved_at TIMESTAMPTZ,
    classification TEXT NOT NULL CHECK (classification IN ('primary', 'secondary', 'tertiary')),
    independence_group TEXT
);

CREATE TABLE events (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    temporal_value JSONB NOT NULL,
    attributes JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE event_participants (
    event_id TEXT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    entity_id TEXT NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
    role TEXT NOT NULL DEFAULT 'participant',
    PRIMARY KEY (event_id, entity_id, role)
);

CREATE TABLE relationships (
    id TEXT PRIMARY KEY,
    from_entity_id TEXT NOT NULL REFERENCES entities(id),
    to_entity_id TEXT NOT NULL REFERENCES entities(id),
    kind TEXT NOT NULL,
    valid_from JSONB,
    valid_to JSONB,
    attributes JSONB NOT NULL DEFAULT '{}'::jsonb,
    CHECK (from_entity_id <> to_entity_id OR kind = 'self_reference')
);

CREATE TABLE claims (
    id TEXT PRIMARY KEY,
    text TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('corroborated', 'supported_limited', 'disputed', 'unverified', 'rejected_superseded')),
    event_id TEXT REFERENCES events(id) ON DELETE SET NULL,
    attributes JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE TABLE evidence (
    id TEXT PRIMARY KEY,
    claim_id TEXT NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
    source_id TEXT NOT NULL REFERENCES sources(id),
    relation TEXT NOT NULL CHECK (relation IN ('supporting', 'contradicting', 'partially_supporting', 'contextual', 'inconclusive')),
    excerpt TEXT NOT NULL,
    location TEXT,
    retrieved_at TIMESTAMPTZ
);

CREATE TABLE conflicts (
    id TEXT PRIMARY KEY,
    subject TEXT NOT NULL,
    disagreement TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE conflict_claims (
    conflict_id TEXT NOT NULL REFERENCES conflicts(id) ON DELETE CASCADE,
    claim_id TEXT NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
    PRIMARY KEY (conflict_id, claim_id)
);

CREATE TABLE research_runs (
    id TEXT PRIMARY KEY,
    subject_id TEXT NOT NULL REFERENCES entities(id),
    status TEXT NOT NULL,
    config JSONB NOT NULL DEFAULT '{}'::jsonb,
    model TEXT,
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ,
    resumed_from TEXT REFERENCES research_runs(id)
);

CREATE INDEX events_temporal_idx ON events USING GIN (temporal_value);
CREATE INDEX evidence_claim_idx ON evidence (claim_id);
CREATE INDEX relationships_entity_idx ON relationships (from_entity_id, to_entity_id);
