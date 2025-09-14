-- name: AdminListUsers :many
SELECT u.id, u.email, u.email_validated,
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
SELECT u.id, u.email, u.email_validated,
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

-- name: AdminCountUsers :one
SELECT COUNT(*)
FROM users u
JOIN profiles p ON u.id = p.user_id
WHERE
  (sqlc.narg('email')::text IS NULL OR u.email ILIKE '%' || sqlc.narg('email')::text || '%')
  AND (sqlc.narg('name')::text IS NULL OR p.name ILIKE '%' || sqlc.narg('name')::text || '%')
  AND (sqlc.narg('email_validated')::boolean IS NULL OR u.email_validated = sqlc.narg('email_validated')::boolean);

-- name: AdminListActivities :many
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
ORDER BY a.id
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
