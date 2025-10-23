-- name: CreatePushToken :one
INSERT INTO push_tokens (user_id, token, platform)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, token) 
DO UPDATE SET 
    platform = EXCLUDED.platform,
    updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: GetPushTokensByUser :many
SELECT * FROM push_tokens 
WHERE user_id = $1;

-- name: GetPushTokensByToken :one
SELECT * FROM push_tokens 
WHERE token = $1;

-- name: DeletePushToken :exec
DELETE FROM push_tokens 
WHERE token = $1;

-- name: DeletePushTokensByUser :exec
DELETE FROM push_tokens 
WHERE user_id = $1;

-- name: GetPushTokensForGroup :many
SELECT pt.* FROM push_tokens pt
JOIN group_members gm ON pt.user_id = gm.user_id
WHERE gm.group_id = $1;
