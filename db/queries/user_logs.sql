-- name: CreateUserLog :exec
INSERT INTO user_logs (id, user_id, action, created_at)
VALUES ($1, $2, $3, $4);

-- name: GetUserLogsByUserID :many
SELECT id, user_id, action, created_at
FROM user_logs
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUserLogsByUserID :one
SELECT COUNT(*) FROM user_logs WHERE user_id = $1;
