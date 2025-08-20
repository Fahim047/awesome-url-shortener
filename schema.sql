CREATE TABLE IF NOT EXISTS url_mappings (
    id BIGSERIAL PRIMARY KEY,
    short_key VARCHAR(10) UNIQUE NOT NULL,
    long_url TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expire_at TIMESTAMPTZ NULL,
    click_count BIGINT DEFAULT 0,
    user_id BIGINT NULL
);

CREATE INDEX IF NOT EXISTS idx_short_key ON url_mappings(short_key);
