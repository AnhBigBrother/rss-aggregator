-- name: CreateUser :one
INSERT INTO
  users (id, NAME, created_at, updated_at, api_key)
VALUES
  (
    $1,
    $2,
    $3,
    $4,
    ENCODE(SHA256(RANDOM()::TEXT::bytea), 'hex')
  )
RETURNING
  *;


-- name: GetUserByApiKey :one
SELECT
  *
FROM
  users
WHERE
  api_key = $1;