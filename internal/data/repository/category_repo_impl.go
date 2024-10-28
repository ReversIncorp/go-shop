package repository

import (
	"database/sql"
	"errors"
	"marketplace/internal/domain/entities"
	repository2 "marketplace/internal/domain/repository"
	"time"
)

type categoryRepositoryImpl struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) repository2.CategoryRepository {
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

func (r *categoryRepositoryImpl) IsBelongsToStore(categoryID, storeID uint64) (bool, error) {
	var existingStoreID uint64
	err := r.db.QueryRow(
		"SELECT store_id FROM categories WHERE id = $1",
		categoryID,
	).Scan(&existingStoreID)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("category not found")
		}
		return false, err
	}

	if existingStoreID == storeID {
		return true, nil
	}

	return false, errors.New("category does not belong to this store")
}

func (r *categoryRepositoryImpl) Save(category entities.Category) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = category.CreatedAt

	_, err := r.db.Exec(`INSERT INTO categories (
                        name,
                        description,
                        store_id, 
                        created_at, 
                        updated_at) 
                        VALUES ($1, $2, $3, $4, $5)`,
		category.Name,
		category.Description,
		category.StoreID,
		category.CreatedAt,
		category.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepositoryImpl) FindByID(id uint64) (entities.Category, error) {
	var category entities.Category
	err := r.db.QueryRow(`SELECT 
    id, 
    name,
    description,
    store_id,
    created_at, 
    updated_at 
	                      FROM categories WHERE id = $1`, id).
		Scan(&category.ID,
			&category.Name,
			&category.Description,
			&category.StoreID,
			&category.CreatedAt,
			&category.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return entities.Category{}, errors.New("category not found")
	}
	if err != nil {
		return entities.Category{}, err
	}

	return category, nil
}

func (r *categoryRepositoryImpl) Update(category entities.Category) error {
	category.UpdatedAt = time.Now()

	//Разрешаем менять все данные кроме айди стора, его нельзя
	_, err := r.db.Exec(`UPDATE categories 
                        SET name = $1, 
                            description = $2, 
                            updated_at = $3
                        WHERE id = $4`,
		category.Name,
		category.Description,
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
    description,
    store_id, 
    created_at, 
    updated_at 
                             FROM categories 
                             WHERE store_id = $1`, storeID)

	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var categories []entities.Category
	for rows.Next() {
		var category entities.Category
		if err := rows.Scan(&category.ID,
			&category.Name,
			&category.Description,
			&category.StoreID,
			&category.CreatedAt,
			&category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
