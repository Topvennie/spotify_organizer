-- +goose Up
-- +goose StatementBegin
ALTER TABLE albums
ADD COLUMN cover_url TEXT;

ALTER TABLE albums
ADD COLUMN cover_id TEXT;

ALTER TABLE shows
ADD COLUMN cover_url TEXT;

ALTER TABLE shows
ADD COLUMN cover_id TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE shows
DROP COLUMN cover_id;

ALTER TABLE shows
DROP COLUMN cover_url;

ALTER TABLE albums
DROP COLUMN cover_id;

ALTER TABLE albums
DROP COLUMN cover_url;
-- +goose StatementEnd
