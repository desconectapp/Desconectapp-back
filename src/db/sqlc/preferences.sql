-- name: GetUserPreferences :many
SELECT a.id, a.name FROM activities a
JOIN users_preference up ON a.id = up.activity_id
WHERE up.user_id = $1
LIMIT $2 OFFSET $3;


-- name: AddPreference :exec
INSERT INTO users_preference (
  user_id, activity_id
) VALUES (
  $1, $2
);

-- name: DeletePreference :exec 
DELETE FROM users_preference
WHERE user_id = $1 AND activity_id = $2;

-- name: BatchAddPreferences :exec
INSERT INTO users_preference (user_id, activity_id)
SELECT sqlc.arg(user_id), id
FROM activities
WHERE id = ANY(sqlc.arg(activity_ids)::int[])
ON CONFLICT DO NOTHING;
