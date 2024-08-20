CREATE TABLE "verify_emails" (
  "id" bigserial PRIMARY KEY,
  "account_id" uuid UNIQUE NOT NULL,
  "secret_code" varchar NOT NULL,
  "is_used" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT 'now()',
  "expired_at" timestamptz NOT NULL DEFAULT (now() + interval '30 minutes')
);

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("account_id");

ALTER TABLE "accounts" ADD COLUMN IF NOT EXISTS "is_email_verified" boolean NOT NULL DEFAULT false;
