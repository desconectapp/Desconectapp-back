-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (
  name, email, password
) VALUES (
  $1, $2, 'password'
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
  set name = $2,
  email = $3,
  age = $4,
  city = $5,
  current_situation = $6,
  profile_complete = true
WHERE id = $1
RETURNING id;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING id;

-- name: CreateActivityRequest :one
INSERT INTO activity_requests (
  user_id, activity_id, description,
  day_of_week, 
  participants_needed
) VALUES (
  $1, $2, $3,
  $4, $5
)
RETURNING *;

-- name: ListActivityRequests :many
SELECT * FROM activity_requests
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetActivities :many
SELECT * FROM activities
ORDER BY category, id DESC
LIMIT $1 OFFSET $2;

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
