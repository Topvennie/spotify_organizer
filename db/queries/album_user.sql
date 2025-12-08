-- name: AlbumUserCreate :one
INSERT INTO album_users (user_id, album_id)
VALUES ($1, $2)
RETURNING id;

-- name: AlbumUserDeleteByUserAlbum :exec
DELETE FROM album_users
WHERE user_id = $1 AND album_id = $2;
