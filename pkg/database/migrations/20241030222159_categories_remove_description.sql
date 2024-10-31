-- +goose Up
ALTER TABLE categories
    DROP COLUMN IF EXISTS description;

-- +goose Down
ALTER TABLE categories
    ADD COLUMN description TEXT NOT NULL DEFAULT '';