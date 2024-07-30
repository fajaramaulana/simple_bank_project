-- name: CreateUser :one
INSERT INTO users ( username, hashed_password, full_name, email ) VALUES ( $1, $2, $3, $4 ) RETURNING user_uuid, username, full_name, email, role;

-- name: GetUserByUserUUID :one
SELECT  user_uuid
       ,username
       ,full_name
       ,email
       ,role
       ,created_at
       ,updated_at
       ,deleted_at
FROM users
WHERE deleted_at = '0001-01-01 00:00:00+00'
AND user_uuid = $1
LIMIT 1;

-- name: GetUserByUsername :one
SELECT  user_uuid
       ,username
       ,full_name
       ,email
       ,role
       ,created_at
       ,updated_at
       ,deleted_at
FROM users
WHERE deleted_at = '0001-01-01 00:00:00+00'
AND username = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT  user_uuid
       ,username
       ,full_name
       ,email
       ,role
       ,created_at
       ,updated_at
       ,deleted_at
FROM users
WHERE deleted_at = '0001-01-01 00:00:00+00'
AND email = $1
LIMIT 1;

-- name: GetDetailLoginByUsername :one
SELECT hashed_password, user_uuid, role, username, email, full_name
FROM users
WHERE deleted_at = '0001-01-01 00:00:00+00'
AND username = $1
LIMIT 1;

-- name: UpdateUser :one
 UPDATE users
SET hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password), password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at), full_name = COALESCE(sqlc.narg(full_name), full_name), email = COALESCE(sqlc.narg(email), email)
WHERE user_uuid = sqlc.arg(user_uuid) RETURNING user_uuid, username, full_name, email, role, created_at, updated_at, deleted_at;

-- name: UpdateUserPassword :one
 UPDATE users
SET hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password), password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at)
WHERE user_uuid = sqlc.arg(user_uuid) RETURNING user_uuid, username, full_name, email, role, created_at, updated_at, deleted_at;


-- name: UpdateUserVerificationEmail :one
UPDATE users
SET verification_email_code = $1, verification_email_expired_at = $2 
WHERE user_uuid = $3 RETURNING user_uuid, verification_email_code, verification_email_expired_at;

-- name: GetUserByVerificationEmailCode :one
SELECT user_uuid, verification_email_code, verification_email_expired_at FROM users WHERE  verification_email_code = $1 AND verified_email_at = '0001-01-01 07:00:00.000 +0700' LIMIT 1;
