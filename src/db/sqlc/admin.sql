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
