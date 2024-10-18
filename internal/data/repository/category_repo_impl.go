package repository

import (
	"database/sql"
	"errors"
	"marketplace/internal/domain/entities"
	repository2 "marketplace/internal/domain/repository"
	"marketplace/pkg/utils"
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

func (r *categoryRepositoryImpl) Save(category entities.Category, uid int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storeExists, err := utils.CheckStoreExists(r.db, category.StoreID)
	if err != nil || !storeExists {
		return errors.New("store not found")
	}

	isOwner, err := utils.CheckUserOwnsStore(r.db, uid, category.StoreID)
	if err != nil || !isOwner {
		return errors.New("user does not own this store")
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

func (r *categoryRepositoryImpl) FindByID(id int64) (entities.Category, error) {
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

func (r *categoryRepositoryImpl) Update(category entities.Category, uid int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var existingCategoryStoreID int64
	err := r.db.QueryRow("SELECT store_id FROM categories WHERE id = $1", category.ID).Scan(&existingCategoryStoreID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("category not found")
		}
		return err
	}

	isOwner, err := utils.CheckUserOwnsStore(r.db, uid, existingCategoryStoreID)
	if err != nil || !isOwner {
		return errors.New("user does not own this store")
	}

	category.UpdatedAt = time.Now()

	//Разрешаем менять все данные кроме айди стора, его нельзя
	_, err = r.db.Exec(`UPDATE categories 
                        SET name = $1, description = $2, updated_at = $3
                        WHERE id = $4`,
		category.Name, category.Description, category.UpdatedAt, category.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepositoryImpl) Delete(id int64, uid int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var existingCategoryStoreID int64
	err := r.db.QueryRow("SELECT store_id FROM categories WHERE id = $1", id).Scan(&existingCategoryStoreID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("category not found")
		}
		return err
	}

	isOwner, err := utils.CheckUserOwnsStore(r.db, uid, existingCategoryStoreID)
	if err != nil || !isOwner {
		return errors.New("user does not own this store")
	}

	_, err = r.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepositoryImpl) FindAllByStore(storeID int64) ([]entities.Category, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.Query(`SELECT id, name, description, store_id, created_at, updated_at 
                             FROM categories 
                             WHERE store_id = $1`, storeID)
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
