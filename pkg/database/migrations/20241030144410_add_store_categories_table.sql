-- +goose Up
CREATE TABLE IF NOT EXISTS store_categories (
                                                store_id BIGINT NOT NULL,
                                                category_id BIGINT NOT NULL,
                                                PRIMARY KEY (store_id, category_id),
                                                FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE,
                                                FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS store_categories;