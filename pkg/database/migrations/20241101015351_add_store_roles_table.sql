-- +goose Up
CREATE TABLE IF NOT EXISTS store_roles (
                                           store_id BIGINT NOT NULL,
                                           user_id BIGINT NOT NULL,
                                           is_owner BOOLEAN NOT NULL DEFAULT FALSE,
                                           PRIMARY KEY (store_id, user_id),
                                           FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE,
                                           FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS store_roles;