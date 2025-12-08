-- +goose Up
-- +goose StatementBegin
ALTER TABLE playlists
ADD COLUMN deleted_at TIMESTAMPTZ DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE playlists
DROP COLUMN deleted_at;
-- +goose StatementEnd
