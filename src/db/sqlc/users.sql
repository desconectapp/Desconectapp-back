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

-- name: CreateProfile :one
UPDATE users
  SET name = $2,
  age = $3,
  city = $4,
  current_situation = $5,
  gender = $6,
  profile_complete = true
WHERE id = $1
RETURNING id;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING id;
