-- name: DirectoryGetByUser :many
SELECT *
FROM directories d
WHERE d.user_id = $1;

-- name: DirectoryCreate :one
INSERT INTO directories (user_id, name, parent_id)
VALUES ($1, $2, $3)
RETURNING id;

-- name: DirectoryUpdate :exec
UPDATE directories
SET name = $2, parent_id = $3
WHERE id = $1;

-- name: DirectoryDeleteByUser :exec
DELETE FROM directories
WHERE user_id = $1;
