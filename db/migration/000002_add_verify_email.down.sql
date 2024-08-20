DROP TABLE IF EXISTS "verify_emails";

ALTER TABLE "accounts" DROP COLUMN IF EXISTS "is_email_verified";
