-- +goose Up
CREATE TABLE IF NOT EXISTS stores (
                                      id BIGSERIAL PRIMARY KEY,
                                      name VARCHAR(100) NOT NULL,
                                      description TEXT NOT NULL,
                                      owner_id BIGINT,
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS stores;
