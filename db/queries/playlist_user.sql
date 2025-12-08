-- name: PlaylistUserCreate :one
INSERT INTO playlist_users (user_id, playlist_id)
VALUES ($1, $2)
RETURNING id;

-- name: PlaylistUserDeleteByUserPlaylist :exec
DELETE FROM playlist_users
WHERE user_id = $1 AND playlist_id = $2;
