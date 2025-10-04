-- name: CreateActivityRequest :one
INSERT INTO activity_requests (
  user_id, activity_id, description,
  week_timeslots, participants_needed,
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
WHERE (sqlc.narg('search')::text IS NULL OR name ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY category, id DESC
LIMIT $1 OFFSET $2;

-- name: GetActivityRequestByUserAndActivityID :one
SELECT * FROM activity_requests
WHERE user_id = $1 AND activity_id = $2;

-- name: DeleteActivityRequest :exec
DELETE FROM activity_requests
WHERE id = $1;

-- name: GetActivityByName :one
SELECT * FROM activities
WHERE name = $1;

-- name: GetActivityByID :one
SELECT * FROM activities
WHERE id = $1;

-- name: CreateActivity :one
INSERT INTO activities (name, icon, category)
VALUES ($1, $2, $3)
RETURNING *;
