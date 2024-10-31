package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
	"time"
)

type storeRepositoryImpl struct {
	db *sql.DB
}

func NewStoreRepository(db *sql.DB) repository.StoreRepository {
	return &storeRepositoryImpl{
		db: db,
	}
}

func (r *storeRepositoryImpl) IsExist(id uint64) (bool, error) {
	var existingStoreID int64
	err := r.db.QueryRow(
		"SELECT id FROM stores WHERE id = $1",
		id,
	).Scan(&existingStoreID)

	if errors.Is(err, sql.ErrNoRows) {
		return false, errors.New("store not found")
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *storeRepositoryImpl) Save(store entities.Store) (uint64, error) {
	store.UpdatedAt = time.Now()
	store.CreatedAt = time.Now()

	var newStoreID uint64
	err := r.db.QueryRow(`INSERT INTO stores 
    (name, 
     description, 
     owner_id, 
     created_at, 
     updated_at) 
                         VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		store.Name,
		store.Description,
		store.OwnerID,
		store.CreatedAt,
		store.UpdatedAt).Scan(&newStoreID)

	if err != nil {
		return 0, err
	}
	return newStoreID, nil
}

func (r *storeRepositoryImpl) FindByID(id uint64) (entities.Store, error) {
	var store entities.Store
	err := r.db.QueryRow(`SELECT 
    id,
    name, 
    description,
    owner_id, 
    created_at, 
    updated_at 
                          FROM stores WHERE id = $1`, id).
		Scan(&store.ID,
			&store.Name,
			&store.Description,
			&store.OwnerID,
			&store.CreatedAt,
			&store.UpdatedAt)

	if err == sql.ErrNoRows {
		return entities.Store{}, errors.New("store not found")
	}
	if err != nil {
		return entities.Store{}, err
	}

	return store, nil
}

func (r *storeRepositoryImpl) Update(store entities.Store) error {
	store.UpdatedAt = time.Now()
	_, err := r.db.Exec(`UPDATE stores SET
                  name = $1,
                  description = $2,
                  updated_at = $3 
              WHERE id = $4`,
		store.Name,
		store.Description,
		store.UpdatedAt,
		store.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *storeRepositoryImpl) Delete(id uint64) error {
	_, err := r.db.Exec("DELETE FROM stores WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *storeRepositoryImpl) FindAll() ([]entities.Store, error) {
	rows, err := r.db.Query(`SELECT 
    id, 
    name,
    description,
    owner_id, 
    created_at, 
    updated_at 
	FROM stores`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stores []entities.Store
	for rows.Next() {
		var store entities.Store
		if err := rows.Scan(&store.ID,
			&store.Name,
			&store.Description,
			&store.OwnerID,
			&store.CreatedAt,
			&store.UpdatedAt); err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stores, nil
}

func (r *storeRepositoryImpl) AttachCategory(storeID, categoryID uint64) error {
	_, err := r.db.Exec(`INSERT INTO store_categories (store_id, category_id) VALUES ($1, $2)`, storeID, categoryID)
	if err != nil {
		return fmt.Errorf("failed to add category to store: %w", err)
	}
	return nil
}

func (r *storeRepositoryImpl) IsCategoryAttached(storeID, categoryID uint64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM store_categories WHERE store_id = $1 AND category_id = $2)`

	err := r.db.QueryRow(query, storeID, categoryID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *storeRepositoryImpl) DetachCategory(storeID, categoryID uint64) error {
	_, err := r.db.Exec(`DELETE FROM store_categories WHERE store_id = $1 AND category_id = $2`, storeID, categoryID)
	if err != nil {
		return err
	}
	return nil
}
