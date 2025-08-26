-- name: GetUser :one
SELECT * FROM profiles
WHERE user_id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM profiles
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (
  email, password
) VALUES (
  $1, $2
)
RETURNING id, email;

-- name: UpdateUser :one
UPDATE users
  SET email = $2
WHERE id = $1
RETURNING id;

-- name: CreateProfile :one
INSERT INTO profiles (
    user_id, name, age, city, current_situation, gender, profile_complete
) VALUES (
    $1, $2, $3, $4, $5, $6, true
)
ON CONFLICT (user_id) DO UPDATE
SET name = EXCLUDED.name,
    age = EXCLUDED.age,
    city = EXCLUDED.city,
    current_situation = EXCLUDED.current_situation,
    gender = EXCLUDED.gender,
    profile_complete = true
RETURNING user_id, name, age, city, current_situation, gender;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING id;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;