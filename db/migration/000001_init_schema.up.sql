CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "account_id" uuid UNIQUE NOT NULL,
  "email" varchar(255) UNIQUE NOT NULL,
  "hash_password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()',
  "updated_at" timestamptz NOT NULL DEFAULT 'now()',
  "deleted_at" timestamptz
);

CREATE TABLE "profiles" (
  "id" bigserial PRIMARY KEY,
  "account_id" uuid NOT NULL,
  "phone_number" varchar(20) NOT NULL,
  "address" varchar,
  "birthday" timestamptz NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()',
  "updated_at" timestamptz NOT NULL DEFAULT 'now()',
  "deleted_at" timestamptz
);

CREATE TABLE "sessions" (
  "id" bigserial PRIMARY KEY,
  "account_id" uuid NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_id" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expired_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "devices" (
  "id" bigserial PRIMARY KEY,
  "device_id" varchar(255) NOT NULL,
  "account_id" uuid NOT NULL,
  "exp_token_device" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()',
  "deleted_at" timestamptz
);

ALTER TABLE "profiles" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("account_id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("account_id");

ALTER TABLE "devices" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("account_id");

CREATE INDEX "index_account_id" ON "accounts" ("account_id"); 
