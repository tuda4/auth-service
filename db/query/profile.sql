-- name: CreateProfile :one
INSERT INTO profiles (account_id, phone_number, birthday, first_name, last_name, address)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateProfile :one
UPDATE profiles
SET phone_number = $2,
    birthday = $3,
    first_name = $4,
    last_name = $5,
    address = COALESCE(sqlc.narg('address'), address),
    updated_at = 'now()'
WHERE account_id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteProfile :exec
UPDATE profiles
SET deleted_at = 'now()'
WHERE account_id = $1 AND deleted_at IS NULL;

-- name: GetOneProfile :one
SELECT *
FROM profiles
WHERE account_id = $1 AND deleted_at IS NULL
LIMIT 1;
