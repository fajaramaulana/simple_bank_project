ALTER TABLE users ADD COLUMN verification_email_code varchar;
ALTER TABLE users ADD COLUMN verified_email_at timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z');