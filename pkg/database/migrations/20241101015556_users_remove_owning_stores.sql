-- +goose Up
ALTER TABLE users
    DROP COLUMN IF EXISTS owning_stores;

-- +goose Down
ALTER TABLE users
    ADD COLUMN owning_stores BIGINT[] DEFAULT '{}';