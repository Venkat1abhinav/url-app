CREATE TABLE IF NOT EXISTS urls (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    link TEXT NOT NULL,
    hash VARCHAR(32) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS urls_hash_idx ON urls (hash);
CREATE INDEX IF NOT EXISTS urls_expires_at_idx ON urls (expires_at) WHERE expires_at IS NOT NULL;
