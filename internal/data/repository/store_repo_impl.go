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

type storeRepositoryImpl struct {
	db *sql.DB
	mu sync.Mutex
}

func NewStoreRepository(db *sql.DB) repository2.StoreRepository {
	return &storeRepositoryImpl{
		db: db,
	}
}

func (r *storeRepositoryImpl) Save(store entities.Store, userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	store.CreatedAt = time.Now()
	store.UpdatedAt = store.CreatedAt

	var newStoreID uint64
	err := r.db.QueryRow(`INSERT INTO stores (name, description, owner_id, created_at, updated_at) 
                         VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		store.Name, store.Description, userID, store.CreatedAt, store.UpdatedAt).Scan(&newStoreID)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`UPDATE users SET owning_stores = array_append(owning_stores, $1) WHERE id = $2`, newStoreID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *storeRepositoryImpl) FindByID(id int64) (entities.Store, error) {
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

func (r *storeRepositoryImpl) Update(store entities.Store, userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storeExists, err := utils.CheckStoreExists(r.db, store.ID)
	if err != nil || !storeExists {
		return errors.New("store not found")
	}

	isOwner, err := utils.CheckUserOwnsStore(r.db, userID, store.ID)
	if err != nil || !isOwner {
		return errors.New("user does not own this store")
	}

	store.UpdatedAt = time.Now()
	_, err = r.db.Exec(`UPDATE stores SET name = $1, description = $2, updated_at = $3 WHERE id = $4`,
		store.Name, store.Description, store.UpdatedAt, store.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *storeRepositoryImpl) Delete(id int64, uid int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storeExists, err := utils.CheckStoreExists(r.db, id)
	if err != nil || !storeExists {
		return errors.New("store not found")
	}

	isOwner, err := utils.CheckUserOwnsStore(r.db, uid, id)
	if err != nil || !isOwner {
		return errors.New("user does not own this store")
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
