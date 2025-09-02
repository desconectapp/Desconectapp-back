-- name: GetUser :one
SELECT * FROM profiles
WHERE user_id = $1 LIMIT 1;

-- name: GetUserById :one
SELECT id, email, email_validated FROM users
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
RETURNING id, email;

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

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

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
