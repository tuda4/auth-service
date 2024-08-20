-- name: CreateSession :one
INSERT INTO sessions (account_id, refresh_token, user_agent, client_id, is_blocked, expired_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetOneSession :one
SELECT *
FROM sessions
WHERE refresh_token = $1 AND expired_at > 'now()'
LIMIT 1;
