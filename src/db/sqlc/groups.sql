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
SELECT u.id, p.name FROM users u
	JOIN profiles p ON u.id = p.user_id
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

-- name: ExitGroup :exec
DELETE FROM group_members
WHERE group_id = $1 AND user_id = $2;

