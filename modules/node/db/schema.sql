CREATE TABLE IF NOT EXISTS message (
    id SERIAL PRIMARY KEY,
    sender TEXT NOT NULL,
    content TEXT NOT NULL,
    bucket_id BIGINT NOT NULL
);
CREATE INDEX IF NOT EXISTS bucket_index ON message (bucket_id);
