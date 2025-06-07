-- name: CreateSession :one
INSERT INTO sessions (
    user_id,
    token,
    refresh_token,
    expires_at,
    refresh_expires_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM sessions
WHERE token = $1 LIMIT 1;

-- name: GetSessionByRefreshToken :one
SELECT * FROM sessions
WHERE refresh_token = $1 LIMIT 1;

-- name: UpdateSessionToken :one
UPDATE sessions
SET token = $2,
    expires_at = $3
WHERE id = $1
RETURNING *;

-- name: DeleteSessionByToken :exec
DELETE FROM sessions
WHERE token = $1;

-- name: DeleteSessionByRefreshToken :exec
DELETE FROM sessions
WHERE refresh_token = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1; 