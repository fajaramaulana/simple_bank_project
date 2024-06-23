CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
ALTER TABLE account ADD COLUMN account_uuid UUID DEFAULT uuid_generate_v4() NOT NULL;
ALTER TABLE entries ADD COLUMN entries_uuid UUID DEFAULT uuid_generate_v4() NOT NULL;
ALTER TABLE transaction ADD COLUMN transaction_uuid UUID DEFAULT uuid_generate_v4() NOT NULL;
CREATE INDEX idx_account_uuid ON account(account_uuid);
CREATE INDEX idx_entries_uuid ON entries(entries_uuid);
CREATE INDEX idx_transaction_uuid ON transaction(transaction_uuid);