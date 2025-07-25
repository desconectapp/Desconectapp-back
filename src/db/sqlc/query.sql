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
  $1, $2, $3
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

-- name: CreateGroup :one
INSERT INTO groups (
  name, description, location, activity_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING id;

-- name: ListGroups :many
WITH selected_groups AS (
  SELECT *
  FROM groups
  ORDER BY created_at DESC
  LIMIT $1 OFFSET $2
)
SELECT 
  g.id,
  g.name,
  g.description,
  g.created_at,
  g.location,
  a.name AS activity,
  a.icon,
  COUNT(gm.user_id) AS members_count
FROM selected_groups g
JOIN activities a ON g.activity_id = a.id
LEFT JOIN group_members gm ON g.id = gm.group_id
GROUP BY g.id, g.name, g.description, g.created_at, g.location, a.name, a.icon
ORDER BY g.created_at DESC;

-- name: GetGroup :one
SELECT 
  g.id,
  g.name,
  g.description,
  g.created_at,
  a.name AS activity,
  a.icon,
  g.location
FROM groups g
JOIN activities a ON g.activity_id = a.id
LEFT JOIN group_members gm ON gm.group_id = g.id
LEFT JOIN users u ON gm.user_id = u.id
WHERE g.id = $1
GROUP BY g.id, g.name, g.description, g.created_at, a.name, a.icon, g.location;

-- name: DeleteGroup :one
DELETE FROM groups
WHERE id = $1
RETURNING id;

-- name: AddUserToGroup :exec
INSERT INTO group_members (
  group_id, user_id
) VALUES (
  $1, $2
);

-- name: BatchAddUserToGroup :exec
INSERT INTO group_members (user_id, group_id)
SELECT id, sqlc.arg(group_id)
FROM users
WHERE id = ANY(sqlc.arg(user_ids)::int[])
ON CONFLICT DO NOTHING;

-- name: GetGroupMembers :many
SELECT u.id, u.name FROM users u
JOIN group_members gm ON u.id = gm.user_id
WHERE gm.group_id = $1;

-- name: ListUserGroups :many
WITH user_groups AS (
  SELECT g.*
  FROM groups g
  JOIN group_members gm ON g.id = gm.group_id
  WHERE gm.user_id = $3
  ORDER BY g.created_at DESC
  LIMIT $1 OFFSET $2
)
SELECT 
  g.id,
  g.name,
  g.description,
  g.created_at,
  g.location,
  a.name AS activity,
  a.icon,
  COUNT(gm_all.user_id) AS members_count
FROM user_groups g
JOIN activities a ON g.activity_id = a.id
LEFT JOIN group_members gm_all ON g.id = gm_all.group_id
GROUP BY g.id, g.name, g.description, g.created_at, g.location, a.name, a.icon
ORDER BY g.created_at DESC;