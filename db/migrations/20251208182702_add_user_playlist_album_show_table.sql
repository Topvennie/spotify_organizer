-- +goose Up
-- +goose StatementBegin
CREATE TABLE playlist_users (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  playlist_id INTEGER NOT NULL REFERENCES playlists (id) ON DELETE CASCADE,

  UNIQUE (user_id, playlist_id)
);

CREATE TABLE album_users (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  album_id INTEGER NOT NULL REFERENCES albums (id) ON DELETE CASCADE,

  UNIQUE(user_id, album_id)
);

CREATE TABLE show_users (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  show_id INTEGER NOT NULL REFERENCES shows (id) ON DELETE CASCADE,

  UNIQUE (user_id, show_id)
);

INSERT INTO playlist_users (user_id, playlist_id)
SELECT user_id, id FROM playlists;

ALTER TABLE playlists
DROP COLUMN user_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE playlists
ADD COLUMN user_id INTEGER REFERENCES users (id);

UPDATE playlists p
SET user_id = up.user_id
FROM user_playlists up
WHERE up.playlist_id = p.id;

ALTER TABLE playlists
ALTER COLUMN user_id SET NOT NULL;

DROP TABLE show_users;
DROP TABLE album_users;
DROP TABLE playlist_users;
-- +goose StatementEnd
