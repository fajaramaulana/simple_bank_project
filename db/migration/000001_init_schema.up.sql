CREATE TABLE "account" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "transaction" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" decimal NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE INDEX ON "account" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transaction" ("from_account_id");

CREATE INDEX ON "transaction" ("to_account_id");

CREATE INDEX ON "transaction" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "account" ("id");

ALTER TABLE "transaction" ADD FOREIGN KEY ("from_account_id") REFERENCES "account" ("id");

ALTER TABLE "transaction" ADD FOREIGN KEY ("to_account_id") REFERENCES "account" ("id");