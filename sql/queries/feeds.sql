-- name: CreateFeed :one
INSERT INTO
  feeds (ids, created_at, updated_at, NAME, url, user_id)
VALUES
  ($1, $2, $3, $4, $5, $6)
RETURNING
  *;