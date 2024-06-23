-- name: CreateEntry :one
INSERT INTO entries (
  account_id,
  amount
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE deleted_at IS NULL AND id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
WHERE deleted_at IS NULL AND account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

