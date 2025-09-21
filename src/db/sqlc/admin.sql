-- name: AdminListUsers :many
SELECT u.id, u.email, u.email_validated, u.is_suspended,
       p.name, p.age, p.city, p.current_situation, p.gender, p.profile_complete, p.created_at
FROM users u
JOIN profiles p ON u.id = p.user_id
WHERE
  (sqlc.narg('email')::text IS NULL OR u.email ILIKE '%' || sqlc.narg('email')::text || '%')
  AND (sqlc.narg('name')::text IS NULL OR p.name ILIKE '%' || sqlc.narg('name')::text || '%')
  AND (sqlc.narg('email_validated')::boolean IS NULL OR u.email_validated = sqlc.narg('email_validated')::boolean)
ORDER BY u.id
LIMIT $1 OFFSET $2;

-- name: AdminGetUser :one
SELECT u.id, u.email, u.email_validated, u.is_suspended,
       p.name, p.age, p.city, p.current_situation, p.gender, p.profile_complete, p.created_at
FROM users u
JOIN profiles p ON u.id = p.user_id
WHERE u.id = $1;

-- name: AdminCreateUser :one
INSERT INTO users (email, password, email_validated)
VALUES ($1, $2, $3)
RETURNING id, email, email_validated;

-- name: AdminCreateProfile :one
INSERT INTO profiles (user_id, name, age, city, current_situation, gender)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING user_id, name, age, city, current_situation, gender, profile_complete, created_at;

-- name: AdminUpdateUser :exec
UPDATE users
SET email = $2, email_validated = $3
WHERE id = $1;

-- name: AdminUpdateProfile :exec
UPDATE profiles
SET name = $2, age = $3, city = $4, current_situation = $5, gender = $6, profile_complete = $7
WHERE user_id = $1;

-- name: AdminDeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: AdminSuspendUser :exec
WITH suspend AS (
    UPDATE users
    SET is_suspended = true
    WHERE id = $1
    RETURNING id
)
DELETE FROM group_members
WHERE user_id IN (SELECT id FROM suspend);

-- name: AdminUnsuspendUser :exec
UPDATE users
SET is_suspended = false
WHERE id = $1
RETURNING id;

-- name: AdminCountUsers :one
SELECT COUNT(*)
FROM users u
JOIN profiles p ON u.id = p.user_id
WHERE
  (sqlc.narg('email')::text IS NULL OR u.email ILIKE '%' || sqlc.narg('email')::text || '%')
  AND (sqlc.narg('name')::text IS NULL OR p.name ILIKE '%' || sqlc.narg('name')::text || '%')
  AND (sqlc.narg('email_validated')::boolean IS NULL OR u.email_validated = sqlc.narg('email_validated')::boolean);

-- name: AdminListActivitiesAsc :many
SELECT 
  a.id, 
  a.name, 
  a.icon, 
  a.category, 
  a.created_at,
  (SELECT COUNT(*) FROM groups g WHERE g.activity_id = a.id) AS group_count,
  (SELECT COUNT(*) FROM partial_matches pm WHERE pm.activity_id = a.id) AS partial_match_count,
  (SELECT COUNT(*) FROM activity_requests ar WHERE ar.activity_id = a.id) AS request_count,
  (SELECT COUNT(*) FROM users_preference up WHERE up.activity_id = a.id) AS user_count
FROM activities a
WHERE
  (sqlc.narg('name')::text IS NULL OR a.name ILIKE '%' || sqlc.narg('name')::text || '%')
  AND (sqlc.narg('category')::categories IS NULL OR a.category = sqlc.narg('category')::categories)
ORDER BY
  CASE sqlc.narg('sort_field')
    WHEN 'name' THEN a.name
    WHEN 'category' THEN a.category::text
    WHEN 'created_at' THEN a.created_at::text
    WHEN 'group_count' THEN (SELECT COUNT(*) FROM groups g WHERE g.activity_id = a.id)::text
    WHEN 'partial_match_count' THEN (SELECT COUNT(*) FROM partial_matches pm WHERE pm.activity_id = a.id)::text
    WHEN 'request_count' THEN (SELECT COUNT(*) FROM activity_requests ar WHERE ar.activity_id = a.id)::text
    WHEN 'user_count' THEN (SELECT COUNT(*) FROM users_preference up WHERE up.activity_id = a.id)::text
  END
  , a.id
  ASC
LIMIT $1 OFFSET $2;

-- name: AdminListActivitiesDesc :many
SELECT 
  a.id, 
  a.name, 
  a.icon, 
  a.category, 
  a.created_at,
  (SELECT COUNT(*) FROM groups g WHERE g.activity_id = a.id) AS group_count,
  (SELECT COUNT(*) FROM partial_matches pm WHERE pm.activity_id = a.id) AS partial_match_count,
  (SELECT COUNT(*) FROM activity_requests ar WHERE ar.activity_id = a.id) AS request_count,
  (SELECT COUNT(*) FROM users_preference up WHERE up.activity_id = a.id) AS user_count
FROM activities a
WHERE
  (sqlc.narg('name')::text IS NULL OR a.name ILIKE '%' || sqlc.narg('name')::text || '%')
  AND (sqlc.narg('category')::categories IS NULL OR a.category = sqlc.narg('category')::categories)
ORDER BY
  CASE sqlc.narg('sort_field')
    WHEN 'name' THEN a.name
    WHEN 'category' THEN a.category::text
    WHEN 'created_at' THEN a.created_at::text
    WHEN 'group_count' THEN (SELECT COUNT(*) FROM groups g WHERE g.activity_id = a.id)::text
    WHEN 'partial_match_count' THEN (SELECT COUNT(*) FROM partial_matches pm WHERE pm.activity_id = a.id)::text
    WHEN 'request_count' THEN (SELECT COUNT(*) FROM activity_requests ar WHERE ar.activity_id = a.id)::text
    WHEN 'user_count' THEN (SELECT COUNT(*) FROM users_preference up WHERE up.activity_id = a.id)::text
  END
  DESC
LIMIT $1 OFFSET $2;

-- name: AdminGetActivity :one
SELECT 
  a.id, 
  a.name, 
  a.icon, 
  a.category, 
  a.created_at,
  (SELECT COUNT(*) FROM groups g WHERE g.activity_id = a.id) AS group_count,
  (SELECT COUNT(*) FROM partial_matches pm WHERE pm.activity_id = a.id) AS partial_match_count,
  (SELECT COUNT(*) FROM activity_requests ar WHERE ar.activity_id = a.id) AS request_count,
  (SELECT COUNT(*) FROM users_preference up WHERE up.activity_id = a.id) AS user_count
FROM activities a
WHERE a.id = $1;

-- name: AdminCreateActivity :one
INSERT INTO activities (name, icon, category)
VALUES ($1, $2, $3)
RETURNING id, name, icon, category, created_at;

-- name: AdminUpdateActivity :exec
UPDATE activities
SET name = $2, icon = $3, category = $4
WHERE id = $1;

-- name: AdminDeleteActivity :exec
DELETE FROM activities WHERE id = $1;

-- name: AdminCountActivities :one
SELECT COUNT(*)
FROM activities a
WHERE
  (sqlc.narg('name')::text IS NULL OR a.name ILIKE '%' || sqlc.narg('name')::text || '%')
  AND (sqlc.narg('category')::categories IS NULL OR a.category = sqlc.narg('category')::categories);

-- name: AdminListGroupsAsc :many
SELECT g.id,
       g.name,
       g.description,
       g.location,
       g.activity_id,
       g.created_at,
       COUNT(m.user_id) AS member_count
FROM groups g
LEFT JOIN group_members m ON g.id = m.group_id
GROUP BY g.id
ORDER BY g.name ASC
LIMIT $1 OFFSET $2;

-- name: AdminListGroupsDesc :many
SELECT g.id,
       g.name,
       g.description,
       g.location,
       g.activity_id,
       g.created_at,
       COUNT(m.user_id) AS member_count
FROM groups g
LEFT JOIN group_members m ON g.id = m.group_id
GROUP BY g.id
ORDER BY g.name DESC
LIMIT $1 OFFSET $2;

-- name: AdminGetGroup :one
SELECT g.id,
       g.name,
       g.description,
       g.location,
       g.activity_id,
       g.created_at,
       COUNT(m.user_id) AS member_count
FROM groups g
LEFT JOIN group_members m ON g.id = m.group_id
WHERE g.id = $1
GROUP BY g.id, g.name, g.description, g.location, g.activity_id, g.created_at;

-- name: AdminCreateGroup :one
INSERT INTO groups (name, description, location, activity_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: AdminUpdateGroup :one
UPDATE groups
SET name = $2,
    description = $3,
    location = $4,
    activity_id = $5
WHERE id = $1
RETURNING *;

-- name: AdminDeleteGroup :exec
DELETE FROM groups
WHERE id = $1;

-- name: AdminListGroupMembers :many
SELECT u.id, p.name, u.email
FROM group_members gm
JOIN users u ON gm.user_id = u.id
JOIN profiles p ON gm.user_id = p.user_id
WHERE gm.group_id = $1;

-- name: AdminAddGroupMember :exec
INSERT INTO group_members (group_id, user_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: AdminRemoveGroupMember :exec
DELETE FROM group_members
WHERE group_id = $1 AND user_id = $2;

-- name: AdminCountGroups :one
SELECT COUNT(*) FROM groups;
