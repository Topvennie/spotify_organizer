-- name: PlaylistGet :one
SELECT *
FROM playlists
WHERE id = $1;

-- name: PlaylistGetBySpotify :one
SELECT *
FROM playlists
WHERE spotify_id = $1 AND deleted_at IS NULL;

-- name: PlaylistGetByUserWithOwner :many
SELECT sqlc.embed(p), sqlc.embed(u)
FROM playlists p
LEFT JOIN playlist_users pu ON pu.playlist_id = p.id
LEFT JOIN users u ON u.uid = p.owner_uid
WHERE pu.user_id = $1 AND deleted_at IS NULL
ORDER BY p.name;

-- name: PlaylistCreate :one
INSERT INTO playlists (spotify_id, owner_uid, name, description, public, track_amount, collaborative, cover_id, cover_url)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;

-- name: PlaylistUpdateBySpotify :exec
UPDATE playlists
SET owner_uid = $2, name = $3, description = $4, public = $5, track_amount = $6, collaborative = $7, cover_id = $8, cover_url = $9
WHERE spotify_id = $1;

-- name: PlaylistDelete :exec
UPDATE playlists
SET deleted_at = NOW()
WHERE id = $1;
