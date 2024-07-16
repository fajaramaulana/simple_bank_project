-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  email,
  password,
  balance,
  currency,
  refresh_token
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetAccount :one
SELECT id, owner, email, currency, balance, refresh_token, created_at, account_uuid, updated_at, deleted_at FROM accounts
WHERE deleted_at IS NULL AND id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT id, owner, email,  currency, balance, refresh_token, created_at, account_uuid, updated_at, deleted_at FROM accounts
WHERE deleted_at IS NULL AND id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT id, owner, email, currency, balance, refresh_token, created_at, account_uuid, updated_at, deleted_at FROM accounts
WHERE deleted_at IS NULL AND  owner LIKE  $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2, updated_at = now()
WHERE id = $1
RETURNING id, owner, email, currency, balance, refresh_token, created_at, account_uuid, updated_at, deleted_at;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING id, owner, email, currency, balance, refresh_token, created_at, account_uuid, updated_at, deleted_at;

-- name: SubtractAccountBalance :one
UPDATE accounts
SET balance = balance - sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING id, owner, email, currency, balance, refresh_token, created_at, account_uuid, updated_at, deleted_at;

-- name: SoftDeleteAccount :exec
UPDATE accounts
SET deleted_at = now()
WHERE id = $1
RETURNING id, owner, email, currency, balance, refresh_token, created_at, account_uuid, updated_at, deleted_at;