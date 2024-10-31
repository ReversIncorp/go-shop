-- +goose Up
ALTER TABLE categories
    DROP COLUMN IF EXISTS store_id;

-- +goose Down
ALTER TABLE categories
    ADD COLUMN store_id BIGINT;