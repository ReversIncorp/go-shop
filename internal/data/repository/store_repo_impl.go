package repository

import (
	"database/sql"
	"errors"
	"marketplace/internal/domain/entities"
	repository2 "marketplace/internal/domain/repository"
	"sync"
	"time"
)

type storeRepositoryImpl struct {
	db *sql.DB
	mu sync.Mutex
}

func NewStoreRepository(db *sql.DB) repository2.StoreRepository {
	return &storeRepositoryImpl{
		db: db,
	}
}

func (r *storeRepositoryImpl) Save(store entities.Store) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var existingStoreID uint64
	err := r.db.QueryRow("SELECT id FROM stores WHERE id = $1", store.ID).Scan(&existingStoreID)
	if err == nil {
		return errors.New("store already exists in the database")
	}
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	store.CreatedAt = time.Now()
	store.UpdatedAt = store.CreatedAt

	_, err = r.db.Exec(`INSERT INTO stores (name, description, owner_id, created_at, updated_at) 
                        VALUES ($1, $2, $3, $4, $5)`,
		store.Name, store.Description, store.OwnerID, store.CreatedAt, store.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *storeRepositoryImpl) FindByID(id uint64) (entities.Store, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var store entities.Store
	err := r.db.QueryRow(`SELECT id, name, description, owner_id, created_at, updated_at 
                          FROM stores WHERE id = $1`, id).
		Scan(&store.ID, &store.Name, &store.Description, &store.OwnerID, &store.CreatedAt, &store.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Store{}, errors.New("store not found")
		}
		return entities.Store{}, err
	}

	return store, nil
}

func (r *storeRepositoryImpl) Update(store entities.Store) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var existingStoreID uint64
	err := r.db.QueryRow("SELECT id FROM stores WHERE id = $1", store.ID).Scan(&existingStoreID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("store not found")
		}
		return err
	}

	store.UpdatedAt = time.Now()

	_, err = r.db.Exec(`UPDATE stores 
                        SET name = $1, description = $2, owner_id = $3, updated_at = $4 
                        WHERE id = $5`,
		store.Name, store.Description, store.OwnerID, store.UpdatedAt, store.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *storeRepositoryImpl) Delete(id uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var existingStoreID uint64
	err := r.db.QueryRow("SELECT id FROM stores WHERE id = $1", id).Scan(&existingStoreID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("store not found")
		}
		return err
	}

	_, err = r.db.Exec("DELETE FROM stores WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (r *storeRepositoryImpl) FindAll() ([]entities.Store, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.Query(`SELECT id, name, description, owner_id, created_at, updated_at FROM stores`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stores []entities.Store
	for rows.Next() {
		var store entities.Store
		if err := rows.Scan(&store.ID, &store.Name, &store.Description, &store.OwnerID, &store.CreatedAt, &store.UpdatedAt); err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stores, nil
}
