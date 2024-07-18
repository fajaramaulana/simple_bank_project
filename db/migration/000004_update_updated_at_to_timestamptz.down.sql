ALTER TABLE "accounts" DROP COLUMN IF EXISTS updated_at;
ALTER TABLE "accounts" DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE "entries" DROP COLUMN IF EXISTS updated_at;
ALTER TABLE "entries" DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE "transactions" DROP COLUMN IF EXISTS updated_at;
ALTER TABLE "transactions" DROP COLUMN IF EXISTS deleted_at;