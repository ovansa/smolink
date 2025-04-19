CREATE TABLE IF NOT EXISTS url_analytics (
    id SERIAL PRIMARY KEY,
    url_id INTEGER REFERENCES urls(id) ON DELETE CASCADE,
    ip_address TEXT,
    user_agent TEXT,
    accessed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
