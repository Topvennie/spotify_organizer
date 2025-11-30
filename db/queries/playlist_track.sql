-- name: PlaylistTrackCreate :one
INSERT INTO playlist_tracks (playlist_id, track_id)
VALUES ($1, $2)
RETURNING id;

-- name: PlaylistTrackDeleteByPlaylistTrack :exec
DELETE FROM playlist_tracks
WHERE playlist_id = $1 AND track_id = $2;
