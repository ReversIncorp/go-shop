-- +goose Up
ALTER TABLE stores
    DROP COLUMN IF EXISTS owner_id;

-- +goose Down
ALTER TABLE stores
    ADD COLUMN owner_id BIGINT;