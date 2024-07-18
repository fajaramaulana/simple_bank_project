CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "user_uuid" UUID UNIQUE DEFAULT uuid_generate_v4() NOT NULL,
  "username" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),
  "role" varchar NOT NULL DEFAULT 'customers',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL NOT NULL DEFAULT('0001-01-01 00:00:00Z'),
  "deleted_at" timestamptz NOT NULL NOT NULL DEFAULT('0001-01-01 00:00:00Z')
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("user_uuid") REFERENCES "users" ("user_uuid");
