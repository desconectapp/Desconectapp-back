-- name: CreateActivityRequest :one
INSERT INTO activity_requests (
  user_id, activity_id, description,
  week_hours, participants_needed,
  maximum_participants, latitude, longitude,
  search_radius
) VALUES (
  $1, $2, $3,
  $4, $5,
  $6, $7, $8,
  $9
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
