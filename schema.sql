CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE monitors (
    id      UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    api_key  TEXT
)