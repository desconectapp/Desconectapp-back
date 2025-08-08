-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (
  email, password
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
  SET email = $2
WHERE id = $1
RETURNING id;

-- name: CreateProfile :one
UPDATE profiles
  SET name = $2,
  age = $3,
  city = $4,
  current_situation = $5,
  gender = $6,
  profile_complete = true
WHERE user_id = $1
RETURNING user_id;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING id;
