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