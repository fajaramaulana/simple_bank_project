ALTER TABLE accounts
ALTER COLUMN updated_at TYPE timestamptz
USING updated_at AT TIME ZONE 'UTC';
ALTER TABLE accounts
ALTER COLUMN deleted_at TYPE timestamptz
USING deleted_at AT TIME ZONE 'UTC';

ALTER TABLE transactions
ALTER COLUMN updated_at TYPE timestamptz
USING updated_at AT TIME ZONE 'UTC';
ALTER TABLE transactions
ALTER COLUMN deleted_at TYPE timestamptz
USING deleted_at AT TIME ZONE 'UTC';

ALTER TABLE entries
ALTER COLUMN updated_at TYPE timestamptz
USING updated_at AT TIME ZONE 'UTC';
ALTER TABLE entries
ALTER COLUMN deleted_at TYPE timestamptz
USING deleted_at AT TIME ZONE 'UTC';