-- name: GetActivities :many
SELECT * FROM activities
ORDER BY category, created_at DESC
LIMIT $1 OFFSET $2;

-- name: CopyActivities :exec
COPY activities(name, emoji, category)
FROM '/db/files/activities.csv'
WITH (FORMAT csv, HEADER true);