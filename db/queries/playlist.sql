-- name: PlaylistGet :one
SELECT *
FROM playlists
WHERE id = $1;

-- name: PlaylistGetByUserWithOwner :many
SELECT sqlc.embed(p), sqlc.embed(u)
FROM playlists p
LEFT JOIN users u ON u.uid = p.owner_uid
WHERE p.user_id = $1
ORDER BY p.name;

-- name: PlaylistCreate :one
INSERT INTO playlists (user_id, spotify_id, owner_uid, name, description, public, track_amount, collaborative, cover_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;

-- name: PlaylistUpdateBySpotify :exec
UPDATE playlists
SET owner_uid = $2, name = $3, description = $4, public = $5, track_amount = $6, collaborative = $7, cover_id = $8
WHERE spotify_id = $1;

-- name: PlaylistDelete :exec
DELETE FROM playlists
WHERE id = $1;
