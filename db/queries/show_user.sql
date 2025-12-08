-- name: ShowUserCreate :one
INSERT INTO show_users (user_id, show_id)
VALUES ($1, $2)
RETURNING id;

-- name: ShowUserDeleteByUserShow :exec
DELETE FROM show_users
WHERE user_id = $1 AND show_id = $2;
