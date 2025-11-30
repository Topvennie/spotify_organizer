-- +goose Up
-- +goose StatementBegin
ALTER TABLE playlists
RENAME COLUMN tracks TO track_amount;

CREATE TABLE tracks (
  id SERIAL PRIMARY KEY,
  spotify_id TEXT NOT NULL,
  name TEXT NOT NULL,
  popularity INTEGER NOT NULL,

  UNIQUE(spotify_id)
);

CREATE TABLE playlist_tracks (
  id SERIAL PRIMARY KEY,
  playlist_id INTEGER NOT NULL REFERENCES playlists (id) ON DELETE CASCADE,
  track_id  INTEGER NOT NULL REFERENCES tracks (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE playlist_tracks;

DROP TABLE tracks;
-- +goose StatementEnd
