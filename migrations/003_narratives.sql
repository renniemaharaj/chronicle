CREATE TABLE narratives (
    id TEXT PRIMARY KEY,
    subject_id TEXT NOT NULL REFERENCES entities(id),
    chapters JSONB NOT NULL CHECK (jsonb_typeof(chapters) = 'array'),
    version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0)
);

ALTER TABLE sources ADD COLUMN metadata JSONB NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE evidence ADD COLUMN metadata JSONB NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE conflicts ADD COLUMN metadata JSONB NOT NULL DEFAULT '{}'::jsonb;
