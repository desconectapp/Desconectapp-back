-- name: GetUser :one
SELECT * FROM profiles
WHERE user_id = $1 LIMIT 1;

-- name: GetUserById :one
SELECT id, uuid, email, email_validated FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM profiles
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (
  email, password
) VALUES (
  $1, $2
)
RETURNING id, uuid, email;

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
SELECT id, uuid, email, password, email_validated FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByUUID :one
SELECT id, uuid, email, email_validated FROM users
WHERE uuid = $1 LIMIT 1;

-- name: CreateEmailVerificationToken :exec
INSERT INTO email_verification_codes  (
	user_id ,code
) VALUES ($1, $2);

-- name: GetVerificationCode :one
SELECT * FROM email_verification_codes
WHERE user_id = $1;

-- name: VerifyEmail :exec
WITH updated AS (
    UPDATE users
    SET email_validated = TRUE
    WHERE id = $1
)
DELETE FROM email_verification_codes
WHERE user_id = $1;

-- name: UpdateVerificationCode :exec
UPDATE email_verification_codes
SET code = $2
WHERE user_id = $1;

-- name: UpdateUserPassword :exec
WITH updated AS (
	UPDATE users
	SET password=$2
	WHERE id = $1
) 
DELETE FROM email_verification_codes
WHERE user_id = $1;
