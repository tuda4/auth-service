-- name: CreateAccount :one
INSERT INTO accounts (account_id, email, hash_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccountByEmail :one
SELECT *
FROM accounts
WHERE email = $1 AND deleted_at IS NULL AND is_email_verified = TRUE
LIMIT 1;

-- name: ChangePassword :exec
UPDATE accounts
SET hash_password = $2, updated_at = 'now()'
WHERE account_id = $1 AND deleted_at IS NULL;

-- name: DeleteAccount :exec
UPDATE accounts
SET deleted_at = 'now()'
WHERE account_id = $1;

-- name: GetProfileAccount :one
SELECT
    a.account_id,
    a.email,
    p.first_name,
    p.last_name,
    p.phone_number,
    p.birthday,
    p.address
FROM accounts as a
LEFT JOIN profiles as p ON p.account_id = a.account_id
WHERE a.email = $1 AND a.deleted_at IS NULL
LIMIT 1;

-- name: UpdateAccountEmail :exec
UPDATE accounts
SET is_email_verified = true, updated_at = 'now()'
WHERE account_id = $1 AND deleted_at IS NULL;
