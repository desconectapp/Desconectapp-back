-- name: CreateGroup :one
WITH inserted_group AS (
  INSERT INTO groups (name, description, location, location_name, activity_id, public, week_timeslots)
  VALUES ($1, $2, $3, $4, $5, $6, $7)
  RETURNING *
), inserted_members AS (
  INSERT INTO group_members (user_id, group_id)
  SELECT u.id, g.id
  FROM users u
  CROSS JOIN inserted_group g
  WHERE u.id = ANY(sqlc.arg(user_ids)::int[])
  ON CONFLICT DO NOTHING
  RETURNING user_id, group_id
)
SELECT 
  g.id,
  g.name,
  g.description,
  g.location_name,
  g.activity_id,
  g.created_at,
  g.public,
  CAST(
    COALESCE(
      array_agg(m.user_id) FILTER (WHERE m.user_id IS NOT NULL),
      ARRAY[]::int[]
    ) AS int[]
  ) AS members,
  a.name AS activity_name,
  a.icon AS activity_icon
FROM inserted_group g
LEFT JOIN inserted_members m ON g.id = m.group_id
JOIN activities a ON g.activity_id = a.id
GROUP BY g.id, g.name, g.description, g.location_name, g.activity_id, g.created_at, a.name, a.icon, g.public;

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
  g.location_name,
  g.avatar_url,
  a.name AS activity,
  a.icon,
  COUNT(gm.user_id) AS members_count
FROM selected_groups g
JOIN activities a ON g.activity_id = a.id
LEFT JOIN group_members gm ON g.id = gm.group_id
GROUP BY g.id, g.name, g.description, g.created_at, g.location_name, a.name, a.icon, g.avatar_url
ORDER BY g.created_at DESC;

-- name: GetGroup :one
SELECT 
  g.id,
  g.name,
  g.description,
  g.created_at,
  a.name AS activity,
  a.icon,
  g.location_name,
  g.public,
  g.avatar_url
FROM groups g
JOIN activities a ON g.activity_id = a.id
LEFT JOIN group_members gm ON gm.group_id = g.id
LEFT JOIN users u ON gm.user_id = u.id
WHERE g.id = $1
GROUP BY g.id, g.name, g.description, g.created_at, a.name, a.icon, g.location_name, g.public;

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
SELECT u.id, u.uuid, p.name, p.avatar_url FROM users u
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
  g.location_name,
  g.avatar_url,
  a.name AS activity,
  a.icon,
  COUNT(gm_all.user_id) AS members_count
FROM user_groups g
JOIN activities a ON g.activity_id = a.id
LEFT JOIN group_members gm_all ON g.id = gm_all.group_id
GROUP BY g.id, g.name, g.description, g.created_at, g.location_name, a.name, a.icon, g.avatar_url
ORDER BY g.created_at DESC;

-- name: ExitGroup :exec
DELETE FROM group_members
WHERE group_id = $1 AND user_id = $2;

-- name: UpdateGroupDescriptiom :exec
UPDATE groups
SET description = $2
WHERE id = $1;

-- name: ChangeGroupPublic :exec
UPDATE groups
SET public = $2
WHERE id = $1;

-- name: ChangeGroupName :exec
UPDATE groups
SET name = $2
WHERE id = $1;

-- name: ChangeGroupLocation :exec
UPDATE groups
SET location = $2, location_name = $3
WHERE id = $1;

-- name: UpdateGroupAvatar :exec
UPDATE groups
SET avatar_url = $2
WHERE id = $1;

-- name: GetOpenGroupsWithFilter :many
SELECT 
    g.id,
    g.name,
    g.description,
    g.location_name,
    g.location,
    g.avatar_url,
    g.week_timeslots,
    a.name AS activity_name,
    a.icon,
    COUNT(gm.user_id) AS member_count
FROM groups g
JOIN activities a ON g.activity_id = a.id
WHERE g.public = true
  AND g.activity_id = sqlc.narg('activity_id')::int
LIMIT $1 OFFSET $2;

-- name: GetOpenGroupsNoFilter :many
SELECT 
    g.id,
    g.name,
    g.description,
    g.location_name,
    g.location,
    g.avatar_url,
    g.week_timeslots,
    a.name AS activity_name,
    a.icon,
    COUNT(gm.user_id) AS member_count
FROM groups g
JOIN activities a ON g.activity_id = a.id
LEFT JOIN group_members gm ON g.id = gm.group_id
WHERE g.public = true
  AND (sqlc.narg('activity_id')::int IS NULL OR g.activity_id = sqlc.narg('activity_id')::int)
GROUP BY g.id, g.name, g.description, g.location_name, g.avatar_url, g.week_timeslots, a.name, a.icon
LIMIT $1 OFFSET $2;

-- name: GetOpenGroupsWithLocation :many
SELECT 
    g.id,
    g.name,
    g.description,
    g.location,
    g.avatar_url,
    g.week_timeslots,
    a.name AS activity_name,
    a.icon,
    COUNT(gm.user_id) AS member_count,
    CAST((6371 * acos(
        cos(radians(sqlc.arg('latitude')::float)) * 
        cos(radians(CAST(split_part(g.location, ',', 1) AS float))) *
        cos(radians(CAST(split_part(g.location, ',', 2) AS float)) - radians(sqlc.arg('longitude')::float)) +
        sin(radians(sqlc.arg('latitude')::float)) * 
        sin(radians(CAST(split_part(g.location, ',', 1) AS float)))
    )) AS float) AS distance_km
FROM groups g
JOIN activities a ON g.activity_id = a.id
LEFT JOIN group_members gm ON g.id = gm.group_id
WHERE g.public = true
  AND g.location IS NOT NULL
  AND g.location != ''
  AND (sqlc.narg('activity_id')::int IS NULL OR g.activity_id = sqlc.narg('activity_id')::int)
  AND (6371 * acos(
        cos(radians(sqlc.arg('latitude')::float)) * 
        cos(radians(CAST(split_part(g.location, ',', 1) AS float))) *
        cos(radians(CAST(split_part(g.location, ',', 2) AS float)) - radians(sqlc.arg('longitude')::float)) +
        sin(radians(sqlc.arg('latitude')::float)) * 
        sin(radians(CAST(split_part(g.location, ',', 1) AS float)))
    )) <= sqlc.arg('radius')::float
GROUP BY g.id, g.name, g.description, g.location, g.avatar_url, g.week_timeslots, a.name, a.icon
ORDER BY distance_km
LIMIT $1 OFFSET $2;

-- name: GetPreferredGroups :many
SELECT g.id,
       g.name,
       g.description,
       g.location_name,
       g.avatar_url,
       g.week_timeslots,
       g.location,
       g.public,
       g.activity_id,
       g.created_at,
       COUNT(gm.user_id) AS member_count,
       a.name   AS activity_name,
       a.icon   AS activity_icon
FROM groups g
JOIN users_preference up ON g.activity_id = up.activity_id
LEFT JOIN group_members gm ON g.id = gm.group_id
JOIN activities a ON g.activity_id = a.id
WHERE up.user_id = $1
  AND g.public = true
  AND NOT EXISTS (
      SELECT 1
      FROM group_members gm2
      WHERE gm2.group_id = g.id
        AND gm2.user_id = $1
  )
GROUP BY g.id, g.name, g.description, g.location_name, g.avatar_url, g.week_timeslots, g.public, g.activity_id, g.created_at,
         a.name, a.icon
ORDER BY member_count ASC, g.created_at DESC
LIMIT $2 OFFSET $3;
