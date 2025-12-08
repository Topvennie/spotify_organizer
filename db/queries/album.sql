-- name: AlbumGetBySpotify :one
SELECT *
FROM albums
WHERE spotify_id = $1;

-- name: AlbumGetByUser :many
SELECT a.*
FROM albums a
LEFT JOIN album_users au on au.album_id = a.id
WHERE au.user_id = $1;

-- name: AlbumCreate :one
INSERT INTO albums (spotify_id, name, track_amount, popularity)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: AlbumUpdate :exec
UPDATE albums
SET name = $2, track_amount = $3, popularity = $4
WHERE id = $1;
