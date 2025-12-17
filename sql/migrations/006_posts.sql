-- +goose Up

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    published_at TIMESTAMP NOT NULL,
    url TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE posts;
