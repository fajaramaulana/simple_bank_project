-- name: CreateTransaction :one
INSERT INTO transactions (
  from_account_id,
  to_account_id,
  amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE deleted_at IS NULL AND id = $1 LIMIT 1;

-- name: ListTransactions :many
SELECT * FROM transactions
WHERE deleted_at IS NULL AND from_account_id = $1 OR to_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

