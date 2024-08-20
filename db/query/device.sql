-- name: CreateDevice :one
INSERT INTO devices (account_id, device_id, exp_token_device)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOneDevice :one
SELECT *
FROM devices
WHERE account_id = $1 AND device_id = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: GetAllDevices :many
SELECT *
FROM devices
WHERE account_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;
