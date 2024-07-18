ALTER TABLE "accounts" DROP COLUMN IF EXISTS account_uuid;
ALTER TABLE "accounts" DROP COLUMN IF EXISTS user_uuid;
ALTER TABLE "entries" DROP COLUMN IF EXISTS entries_uuid;
ALTER TABLE "transactions" DROP COLUMN IF EXISTS transaction_uuid;
DROP INDEX IF EXISTS idx_account_uuid;
DROP INDEX IF EXISTS idx_user_uuid;
DROP INDEX IF EXISTS idx_entries_uuid;
DROP INDEX IF EXISTS idx_transaction_uuid;