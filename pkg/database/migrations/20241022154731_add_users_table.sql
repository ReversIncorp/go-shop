-- +goose Up
CREATE TABLE IF NOT EXISTS users (
                                     id BIGSERIAL PRIMARY KEY,
                                     name VARCHAR(100) NOT NULL,
                                     email VARCHAR(255) NOT NULL UNIQUE,
                                     password VARCHAR(255) NOT NULL,
                                     is_seller BOOLEAN DEFAULT FALSE,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     owning_stores BIGINT[] DEFAULT '{}'
);

-- +goose Down
DROP TABLE IF EXISTS users;
