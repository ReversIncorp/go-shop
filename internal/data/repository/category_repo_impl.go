package repository

import (
	"database/sql"
	"errors"
	"marketplace/internal/domain/entities"
	repository2 "marketplace/internal/domain/repository"
	"sync"
	"time"
)

type categoryRepositoryImpl struct {
	db *sql.DB
	mu sync.Mutex
}

func NewCategoryRepository(db *sql.DB) repository2.CategoryRepository {
	return &categoryRepositoryImpl{
		db: db,
	}
}

func (r *categoryRepositoryImpl) Save(category entities.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var existingCategoryID uint64
	err := r.db.QueryRow("SELECT id FROM categories WHERE id = $1", category.ID).Scan(&existingCategoryID)
	if err == nil {
		return errors.New("category already exists")
	}
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	category.CreatedAt = time.Now()
	category.UpdatedAt = category.CreatedAt

	_, err = r.db.Exec(`INSERT INTO categories (name, description, store_id, created_at, updated_at) 
                        VALUES ($1, $2, $3, $4, $5)`,
		category.Name, category.Description, category.StoreID, category.CreatedAt, category.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepositoryImpl) FindByID(id uint64) (entities.Category, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var category entities.Category
	err := r.db.QueryRow(`SELECT id, name, description, store_id, created_at, updated_at 
	                      FROM categories WHERE id = $1`, id).
		Scan(&category.ID, &category.Name, &category.Description, &category.StoreID, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Category{}, errors.New("category not found")
		}
		return entities.Category{}, err
	}

	return category, nil
}

func (r *categoryRepositoryImpl) Update(category entities.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var existingCategoryID uint64
	err := r.db.QueryRow("SELECT id FROM categories WHERE id = $1", category.ID).Scan(&existingCategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("category not found")
		}
		return err
	}

	category.UpdatedAt = time.Now()

	// Выполняем SQL-запрос на обновление данных категории
	_, err = r.db.Exec(`UPDATE categories 
                        SET name = $1, description = $2, store_id = $3, updated_at = $4 
                        WHERE id = $5`,
		category.Name, category.Description, category.StoreID, category.UpdatedAt, category.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepositoryImpl) Delete(id uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var existingCategoryID uint64
	err := r.db.QueryRow("SELECT id FROM categories WHERE id = $1", id).Scan(&existingCategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("category not found")
		}
		return err
	}

	_, err = r.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepositoryImpl) FindAll() ([]entities.Category, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.Query(`SELECT id, name, description, store_id, created_at, updated_at FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []entities.Category
	for rows.Next() {
		var category entities.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.StoreID, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
