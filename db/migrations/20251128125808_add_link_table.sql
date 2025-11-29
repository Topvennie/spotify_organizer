-- +goose Up
-- +goose StatementBegin
CREATE TABLE links (
  id SERIAL PRIMARY KEY,
  source_directory_id INTEGER REFERENCES directories (id) ON DELETE CASCADE,
  source_playlist_id INTEGER REFERENCES playlists (id) ON DELETE CASCADE,
  target_directory_id INTEGER REFERENCES directories (id) ON DELETE CASCADE,
  target_playlist_id INTEGER REFERENCES playlists (id) ON DELETE CASCADE,

  CONSTRAINT source_exclusive CHECK (
    (source_directory_id IS NOT NULL AND source_playlist_id IS NULL) OR
    (source_directory_id IS NULL AND source_playlist_id IS NOT NULL)
  ),

  CONSTRAINT target_exclusive CHECK (
    (target_directory_id IS NOT NULL AND target_playlist_id IS NULL) OR
    (target_directory_id IS NULL AND target_playlist_id IS NOT NULL)
  )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE links;
-- +goose StatementEnd
