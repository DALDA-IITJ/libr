CREATE TABLE IF NOT EXISTS message (
    id SERIAL PRIMARY KEY,
    sender TEXT NOT NULL,
    content TEXT NOT NULL,
    bucket_id BIGINT NOT NULL,
    INDEX bucket_index (bucket_id)
)