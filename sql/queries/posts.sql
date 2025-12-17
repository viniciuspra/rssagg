-- name: CreatePost :one
INSERT INTO posts (id, feed_id, title, description, published_at, url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
