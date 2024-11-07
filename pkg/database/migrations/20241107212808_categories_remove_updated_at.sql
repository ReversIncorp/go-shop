-- +goose Up
ALTER TABLE categories
    DROP COLUMN IF EXISTS updated_at;

-- +goose Down
ALTER TABLE categories
    ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;