-- name: DirectoryPlaylistGetByDirectory :many
SELECT *
FROM directory_playlists
WHERE directory_id = ANY($1::int[]);

-- name: DirectoryPlaylistCreate :one
INSERT INTO directory_playlists (directory_id, playlist_id)
VALUES ($1, $2)
RETURNING id;

-- name: DirectoryPlaylistDeleteByDirectoryPlaylist :exec
DELETE FROM directory_playlists
WHERE directory_id = $1 AND playlist_id = $2;

