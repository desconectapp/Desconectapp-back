-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (
  name, email
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
  set name = $2,
  email = $3
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
ORDER BY category, created_at DESC
LIMIT $1 OFFSET $2;

-- name: CopyActivities :exec
COPY activities(name, emoji, category)
FROM '/db/files/activities.csv'
WITH (FORMAT csv, HEADER true);