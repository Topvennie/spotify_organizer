-- name: LinkGetByUser :many
SELECT l.*
FROM links l
LEFT JOIN directories d ON d.id = l.source_directory_id
LEFT JOIN playlists p ON p.id = l.source_playlist_id
WHERE d.user_id = $1 OR p.user_id = $1;

-- name: LinkCreate :one
INSERT INTO links (source_directory_id, source_playlist_id, target_directory_id, target_playlist_id)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: LinkUpdate :exec
UPDATE links
SET source_directory_id = $2, source_playlist_id = $3, target_directory_id = $4, target_playlist_id = $5
WHERE id = $1;

-- name: LinkDelete :exec
DELETE FROM links
WHERE id = $1;
