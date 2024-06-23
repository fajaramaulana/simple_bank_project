ALTER TABLE account
ALTER COLUMN updated_at TYPE timestamptz
USING updated_at AT TIME ZONE 'UTC';
ALTER TABLE account
ALTER COLUMN deleted_at TYPE timestamptz
USING deleted_at AT TIME ZONE 'UTC';

ALTER TABLE transaction
ALTER COLUMN updated_at TYPE timestamptz
USING updated_at AT TIME ZONE 'UTC';
ALTER TABLE transaction
ALTER COLUMN deleted_at TYPE timestamptz
USING deleted_at AT TIME ZONE 'UTC';

ALTER TABLE entries
ALTER COLUMN updated_at TYPE timestamptz
USING updated_at AT TIME ZONE 'UTC';
ALTER TABLE entries
ALTER COLUMN deleted_at TYPE timestamptz
USING deleted_at AT TIME ZONE 'UTC';