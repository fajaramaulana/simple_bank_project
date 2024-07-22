// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users ( username, hashed_password, full_name, email ) VALUES ( $1, $2, $3, $4 ) RETURNING user_uuid, username, full_name, email, role
`

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

type CreateUserRow struct {
	UserUuid uuid.UUID `json:"user_uuid"`
	Username string    `json:"username"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email,
	)
	var i CreateUserRow
	err := row.Scan(
		&i.UserUuid,
		&i.Username,
		&i.FullName,
		&i.Email,
		&i.Role,
	)
	return i, err
}

const getDetailLoginByUsername = `-- name: GetDetailLoginByUsername :one
SELECT hashed_password, user_uuid, role, username, email, full_name
FROM users
WHERE deleted_at = '0001-01-01 00:00:00+00'
AND username = $1
LIMIT 1
`

type GetDetailLoginByUsernameRow struct {
	HashedPassword string    `json:"hashed_password"`
	UserUuid       uuid.UUID `json:"user_uuid"`
	Role           string    `json:"role"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	FullName       string    `json:"full_name"`
}

func (q *Queries) GetDetailLoginByUsername(ctx context.Context, username string) (GetDetailLoginByUsernameRow, error) {
	row := q.db.QueryRowContext(ctx, getDetailLoginByUsername, username)
	var i GetDetailLoginByUsernameRow
	err := row.Scan(
		&i.HashedPassword,
		&i.UserUuid,
		&i.Role,
		&i.Username,
		&i.Email,
		&i.FullName,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
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
LIMIT 1
`

type GetUserByEmailRow struct {
	UserUuid  uuid.UUID `json:"user_uuid"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i GetUserByEmailRow
	err := row.Scan(
		&i.UserUuid,
		&i.Username,
		&i.FullName,
		&i.Email,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getUserByUserUUID = `-- name: GetUserByUserUUID :one
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
LIMIT 1
`

type GetUserByUserUUIDRow struct {
	UserUuid  uuid.UUID `json:"user_uuid"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (q *Queries) GetUserByUserUUID(ctx context.Context, userUuid uuid.UUID) (GetUserByUserUUIDRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByUserUUID, userUuid)
	var i GetUserByUserUUIDRow
	err := row.Scan(
		&i.UserUuid,
		&i.Username,
		&i.FullName,
		&i.Email,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
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
LIMIT 1
`

type GetUserByUsernameRow struct {
	UserUuid  uuid.UUID `json:"user_uuid"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (GetUserByUsernameRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i GetUserByUsernameRow
	err := row.Scan(
		&i.UserUuid,
		&i.Username,
		&i.FullName,
		&i.Email,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
 UPDATE users
SET hashed_password = COALESCE($1, hashed_password), password_changed_at = COALESCE($2, password_changed_at), full_name = COALESCE($3, full_name), email = COALESCE($4, email)
WHERE user_uuid = $5 RETURNING user_uuid, username, full_name, email, role, created_at, updated_at, deleted_at
`

type UpdateUserParams struct {
	HashedPassword    sql.NullString `json:"hashed_password"`
	PasswordChangedAt sql.NullTime   `json:"password_changed_at"`
	FullName          sql.NullString `json:"full_name"`
	Email             sql.NullString `json:"email"`
	UserUuid          uuid.UUID      `json:"user_uuid"`
}

type UpdateUserRow struct {
	UserUuid  uuid.UUID `json:"user_uuid"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (UpdateUserRow, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.HashedPassword,
		arg.PasswordChangedAt,
		arg.FullName,
		arg.Email,
		arg.UserUuid,
	)
	var i UpdateUserRow
	err := row.Scan(
		&i.UserUuid,
		&i.Username,
		&i.FullName,
		&i.Email,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :one
 UPDATE users
SET hashed_password = COALESCE($1, hashed_password), password_changed_at = COALESCE($2, password_changed_at)
WHERE user_uuid = $3 RETURNING user_uuid, username, full_name, email, role, created_at, updated_at, deleted_at
`

type UpdateUserPasswordParams struct {
	HashedPassword    sql.NullString `json:"hashed_password"`
	PasswordChangedAt sql.NullTime   `json:"password_changed_at"`
	UserUuid          uuid.UUID      `json:"user_uuid"`
}

type UpdateUserPasswordRow struct {
	UserUuid  uuid.UUID `json:"user_uuid"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (UpdateUserPasswordRow, error) {
	row := q.db.QueryRowContext(ctx, updateUserPassword, arg.HashedPassword, arg.PasswordChangedAt, arg.UserUuid)
	var i UpdateUserPasswordRow
	err := row.Scan(
		&i.UserUuid,
		&i.Username,
		&i.FullName,
		&i.Email,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
