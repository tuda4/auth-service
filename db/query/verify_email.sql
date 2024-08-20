-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
    account_id,
    secret_code,
    is_used
) VALUES (
    $1,
    $2,
    $3
    )
RETURNING *;

-- name: GetVerifyEmailBySecretCode :one
SELECT * FROM verify_emails WHERE account_id = $1 AND secret_code = $2 AND is_used = false AND expired_at > now() LIMIT 1;

-- name: UpdateVerifyEmail :exec
UPDATE verify_emails SET is_used = true WHERE account_id = $1 AND secret_code = $2 AND is_used = false AND expired_at > now();
