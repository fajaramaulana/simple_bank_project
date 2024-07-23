-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balance,
  user_uuid,
  currency,
  status
) VALUES (
  $1, $2, $3, $4, 1
) RETURNING id, owner,  currency, balance, user_uuid, account_uuid, status, created_at;

-- name: GetAccount :one
SELECT id, owner, currency, balance, user_uuid, created_at, account_uuid, updated_at, deleted_at, status FROM accounts
WHERE deleted_at IS NULL AND id = $1 LIMIT 1;

-- name: GetAccountByUUID :one
SELECT id, owner, currency, balance, user_uuid, created_at, account_uuid, updated_at, deleted_at, status FROM accounts
WHERE deleted_at IS NULL AND account_uuid = $1 LIMIT 1;

-- name: GetAccountByUserUUIDAndCurrency :one
SELECT id, owner, currency, balance, user_uuid, created_at, account_uuid, updated_at, deleted_at, status FROM accounts
WHERE deleted_at IS NULL AND user_uuid = $1 AND currency = $2 LIMIT 1;

-- name: GetAccountByUserUUID :one
SELECT id, owner, currency, balance, user_uuid, created_at, account_uuid, updated_at, deleted_at, status FROM accounts
WHERE deleted_at IS NULL AND user_uuid = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT id, owner,  currency, balance, user_uuid, created_at, account_uuid, updated_at, deleted_at, status FROM accounts
WHERE deleted_at IS NULL AND id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT accounts.id, owner, currency, balance, accounts.user_uuid, accounts.created_at, account_uuid, accounts.updated_at, accounts.deleted_at, status, u.email, u.full_name, u.username FROM accounts
LEFT JOIN users u ON accounts.user_uuid = u.user_uuid
WHERE accounts.deleted_at IS NULL
ORDER BY accounts.id
LIMIT $1
OFFSET $2;

-- name: CountAccounts :one
SELECT COUNT(*) FROM accounts
WHERE deleted_at IS NULL;

-- name: ListAccountsByUserUUID :many
SELECT accounts.id, owner, currency, balance, accounts.user_uuid, accounts.created_at, account_uuid, accounts.updated_at, accounts.deleted_at, status, u.email, u.full_name, u.username FROM accounts
LEFT JOIN users u ON accounts.user_uuid = u.user_uuid
WHERE accounts.deleted_at IS NULL AND accounts.user_uuid = $1
ORDER BY accounts.id
LIMIT $2
OFFSET $3;

-- name: CountAccountsByUserUUID :one
SELECT COUNT(*) FROM accounts
WHERE deleted_at IS NULL AND user_uuid = $1;


-- name: UpdateAccount :one
WITH updated_account AS (
    UPDATE accounts
    SET balance = $2, updated_at = now()
    WHERE accounts.id = $1
    RETURNING id, owner, currency, balance, user_uuid, created_at, account_uuid, updated_at, deleted_at, status
)
SELECT
    ua.id AS account_id,
    ua.owner,
    ua.currency,
    ua.balance,
    ua.user_uuid,
    ua.created_at,
    ua.account_uuid,
    ua.updated_at,
    ua.deleted_at,
    ua.status,
    u.full_name,  -- Add desired user columns here
    u.username,  -- Add desired user columns here
    u.email      -- Add desired user columns here
FROM
    updated_account ua
LEFT JOIN
    users u ON ua.user_uuid = u.user_uuid;



-- name: UpdateProfileAccount :one
WITH updated_account AS (
  UPDATE accounts
  SET owner = $2, currency = $3, status = $4, updated_at = now()
  WHERE account_uuid = $1
  RETURNING id, owner, currency, balance, user_uuid, created_at, account_uuid, updated_at, deleted_at, status
)
SELECT
    ua.id AS account_id,
    ua.owner,
    ua.currency,
    ua.balance,
    ua.user_uuid,
    ua.created_at,
    ua.account_uuid,
    ua.updated_at,
    ua.deleted_at,
    ua.status,
    u.full_name,  -- Add desired user columns here
    u.username,  -- Add desired user columns here
    u.email      -- Add desired user columns here
FROM
    updated_account ua
LEFT JOIN
    users u ON ua.user_uuid = u.user_uuid;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING id, owner, currency, balance, user_uuid, created_at, account_uuid, updated_at, deleted_at, status;

-- name: SubtractAccountBalance :one
UPDATE accounts
SET balance = balance - sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING id, owner, currency, balance, user_uuid, created_at, account_uuid, updated_at, deleted_at, status;

-- name: SoftDeleteAccount :exec
UPDATE accounts
SET deleted_at = now()
WHERE id = $1
RETURNING id, owner, currency, balance, user_uuid, created_at, account_uuid, updated_at, deleted_at, status;