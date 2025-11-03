-- name: UserGet :one
SELECT *
FROM users
WHERE id = $1 LIMIT 1;

-- name: UserGetByUID :one
SELECT *
FROM users
WHERE uid = $1 LIMIT 1;

-- name: UserCreate :one
INSERT INTO users (name, email, uid)
VALUES ($1, $2, $3)
RETURNING id;
