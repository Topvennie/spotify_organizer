-- name: ShowGetBySpotify :one
SELECT *
FROM shows
WHERE spotify_id = $1;

-- name: ShowGetByUser :many
SELECT s.*
FROM shows s
LEFT JOIN show_users su on su.show_id = s.id
WHERE su.user_id = $1;

-- name: ShowCreate :one
INSERT INTO shows (spotify_id, episode_amount, name, cover_id, cover_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: ShowUpdate :exec
UPDATE shows
SET name = $2, episode_amount = $3, cover_id = $4, cover_url = $5
WHERE id = $1;
