CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE monitors (
    id      UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    api_key_hash TEXT NOT NULL,
    query TEXT NOT NULL,
    webhook_url TEXT NOT NULL,
    schedule TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    monitor_id UUID NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    title TEXT NOT NULL,
    snippet TEXT NOT NULL,
    found_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON results(monitor_id);