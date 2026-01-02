-- +goose Up
-- +goose StatementBegin
-- Interviews
CREATE TABLE IF NOT EXISTS interviews (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    message_count INTEGER NOT NULL DEFAULT 0 CHECK (message_count >= 0),
    prompt_tokens  INTEGER NOT NULL DEFAULT 0 CHECK (prompt_tokens >= 0),
    completion_tokens  INTEGER NOT NULL DEFAULT 0 CHECK (completion_tokens>= 0),
    cost REAL NOT NULL DEFAULT 0.0 CHECK (cost >= 0.0),
    summary_message_id TEXT,
    todos TEXT,
    difficulty TEXT,
    topic TEXT,
    status TEXT,
    updated_at INTEGER NOT NULL,  -- Timestamp in seconds
    created_at INTEGER NOT NULL   -- Timestamp in seconds
);

CREATE TRIGGER IF NOT EXISTS update_interviews_updated_at
AFTER UPDATE ON interviews
BEGIN
UPDATE interviews SET updated_at = strftime('%s', 'now')
WHERE id = new.id;
END;
