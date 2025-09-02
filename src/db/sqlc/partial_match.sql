
-- name: CreatePartialMatch :one
INSERT INTO partial_matches (
  activity_id, description, week_hours,
  participants_needed, maximum_participants,
  latitude, longitude, members_count,
  search_radius
) VALUES (
  $1, $2, $3,
  $4, $5,
  $6, $7, $8,
  $9
)
RETURNING *;

-- name: FindPartialMatches :many
SELECT * FROM partial_matches
WHERE activity_id = $1
ORDER BY created_at DESC;

-- name: AddUserToPartialMatch :exec
INSERT INTO partial_match_members (
  partial_match_id, user_id
) VALUES (
  $1, $2
) ON CONFLICT DO NOTHING;

-- name: BatchAddUsersToPartialMatch :exec
INSERT INTO partial_match_members (partial_match_id, user_id)
SELECT sqlc.arg(partial_match_id), id
FROM users
WHERE id = ANY(sqlc.arg(user_ids)::int[])
ON CONFLICT DO NOTHING;

-- name: GetPartialMatchesByActivity :many
SELECT pm.*, COUNT(pmm.user_id) as current_members
FROM partial_matches pm
LEFT JOIN partial_match_members pmm ON pm.id = pmm.partial_match_id
WHERE ($1::int IS NULL OR pm.activity_id = $1)
GROUP BY pm.id
ORDER BY pm.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetPartialMatchMembers :many
SELECT u.id, p.name FROM users u
JOIN profiles p ON u.id = p.user_id
JOIN partial_match_members pmm ON u.id = pmm.user_id
WHERE pmm.partial_match_id = $1;

-- name: DeletePartialMatch :one
DELETE FROM partial_matches
WHERE id = $1
RETURNING id;

-- name: DeletePartialMatchesByUserAndActivityID :exec
DELETE FROM partial_matches
WHERE id IN (
  SELECT pm.id
  FROM partial_matches pm
  JOIN partial_match_members pmm ON pm.id = pmm.partial_match_id
  WHERE pmm.user_id = $1 AND pm.activity_id = $2
);