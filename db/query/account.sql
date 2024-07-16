-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  email,
  password,
  balance,
  currency,
  refresh_token,
  status
) VALUES (
  $1, $2, $3, $4, $5, $6, 1
) RETURNING id, owner, email,  currency, balance, refresh_token, account_uuid, status, created_at;

-- name: GetAccount :one
SELECT id, owner, email, currency, balance, refresh_token, created_at, account_uuid, updated_at, deleted_at, status FROM accounts
WHERE deleted_at IS NULL AND id = $1 LIMIT 1;

-- name: GetAccountByUUID :one
SELECT id, owner, email, currency, balance, refresh_token, created_at, account_uuid, updated_at, deleted_at, status FROM accounts
WHERE deleted_at IS NULL AND account_uuid = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT id, owner, email,  currency, balance, refresh_token, created_at, account_uuid, updated_at, deleted_at, status FROM accounts
WHERE deleted_at IS NULL AND id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: GetAccountByEmail :one
SELECT id, account_uuid, owner, email, password, status FROM accounts
WHERE deleted_at IS NULL AND email = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT id, owner, email, currency, balance, refresh_token, created_at, account_uuid, updated_at, deleted_at , status FROM accounts
WHERE deleted_at IS NULL
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CountAccounts :one
SELECT COUNT(*) FROM accounts
WHERE deleted_at IS NULL;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2, updated_at = now()
WHERE id = $1
RETURNING id, owner, email, currency, balance, refresh_token, created_at, account_uuid, updated_at, deleted_at, status;

-- name: UpdateProfileAccount :one
UPDATE accounts
SET owner = $2, currency = $3, status = $4, updated_at = now()
WHERE account_uuid = $1
RETURNING id, owner, email, currency, balance, refresh_token, created_at, account_uuid, updated_at, deleted_at, status;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING id, owner, email, currency, balance, refresh_token, created_at, account_uuid, updated_at, deleted_at, status;

-- name: SubtractAccountBalance :one
UPDATE accounts
SET balance = balance - sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING id, owner, email, currency, balance, refresh_token, created_at, account_uuid, updated_at, deleted_at, status;

-- name: SoftDeleteAccount :exec
UPDATE accounts
SET deleted_at = now()
WHERE id = $1
RETURNING id, owner, email, currency, balance, refresh_token, created_at, account_uuid, updated_at, deleted_at, status;