package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
	"time"
)

type categoryRepositoryImpl struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) repository.CategoryRepository {
	return &categoryRepositoryImpl{
		db: db,
	}
}

func (r *categoryRepositoryImpl) IsExist(id uint64) (bool, error) {
	var existingStoreID uint64
	err := r.db.QueryRow(
		"SELECT id FROM categories WHERE id = $1",
		id,
	).Scan(&existingStoreID)

	if errors.Is(err, sql.ErrNoRows) {
		return false, errors.New("category not found")
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *categoryRepositoryImpl) Save(category entities.Category) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = category.CreatedAt

	_, err := r.db.Exec(`INSERT INTO categories (
                        name,
                        created_at, 
                        updated_at) 
                        VALUES ($1, $2, $3)`,
		category.Name,
		category.CreatedAt,
		category.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepositoryImpl) Update(category entities.Category) error {
	category.UpdatedAt = time.Now()

	_, err := r.db.Exec(`UPDATE categories 
                        SET name = $1, 
                            updated_at = $2
                        WHERE id = $3`,
		category.Name,
		category.UpdatedAt,
		category.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepositoryImpl) Delete(id uint64) error {
	_, err := r.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *categoryRepositoryImpl) FindAllByStore(storeID uint64) ([]entities.Category, error) {
	rows, err := r.db.Query(`SELECT 
		id, 
		name, 
		created_at, 
		updated_at 
	FROM categories 
	INNER JOIN store_categories ON categories.id = store_categories.category_id 
	WHERE store_categories.store_id = $1`, storeID)

	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []entities.Category
	for rows.Next() {
		var category entities.Category
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan category row: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return categories, nil
}
