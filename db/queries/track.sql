-- name: TrackGetBySpotify :one
SELECT *
FROM tracks
WHERE spotify_id = $1;

-- name: TrackGetByPlaylist :many
SELECT t.*
FROM tracks t
LEFT JOIN playlist_tracks p_t ON p_t.track_id = t.id
WHERE p_t.playlist_id = $1;

-- name: TrackCreate :one
INSERT INTO tracks (spotify_id, name, popularity)
VALUES ($1, $2, $3)
RETURNING id;

-- name: TrackUpdateBySpotify :exec
UPDATE tracks
SET name = $2, popularity = $3
WHERE spotify_id = $1;
