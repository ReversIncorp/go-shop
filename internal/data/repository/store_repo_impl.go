package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
	"marketplace/pkg/error_handling"
	"strings"
	"time"

	"github.com/ztrue/tracerr"
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
		return false, tracerr.Wrap(errors.New("store not found"))
	}
	if err != nil {
		return false, tracerr.Wrap(err)
	}

	return true, nil
}

func (r *storeRepositoryImpl) Save(store entities.Store, uid uint64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return tracerr.Wrap(fmt.Errorf("failed to begin transaction: %w", err))
	}

	store.UpdatedAt = time.Now()
	store.CreatedAt = time.Now()

	var newStoreID uint64
	err = tx.QueryRow(`INSERT INTO stores 
		(name, description, created_at, updated_at) 
		VALUES ($1, $2, $3, $4) RETURNING id`,
		store.Name,
		store.Description,
		store.CreatedAt,
		store.UpdatedAt).Scan(&newStoreID)

	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return tracerr.Wrap(err)
		}
		return tracerr.Wrap(fmt.Errorf("failed to save store: %w", err))
	}

	_, err = tx.Exec(`INSERT INTO store_roles (store_id, user_id, is_owner) VALUES ($1, $2, $3)`,
		newStoreID, uid, true)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return err
		}
		return tracerr.Wrap(fmt.Errorf("failed to add user as store admin: %w", err))
	}

	if err = tx.Commit(); err != nil {
		return tracerr.Wrap(fmt.Errorf("failed to commit transaction: %w", err))
	}

	return nil
}

func (r *storeRepositoryImpl) FindByID(id uint64) (entities.Store, error) {
	var store entities.Store
	err := r.db.QueryRow(`SELECT 
    id,
    name, 
    description,
    created_at, 
    updated_at 
                          FROM stores WHERE id = $1`, id).
		Scan(&store.ID,
			&store.Name,
			&store.Description,
			&store.CreatedAt,
			&store.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return entities.Store{}, errorHandling.ErrStoreNotFound
	}
	if err != nil {
		return entities.Store{}, tracerr.Wrap(err)
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
		return tracerr.Wrap(err)
	}

	return nil
}

func (r *storeRepositoryImpl) Delete(id uint64) error {
	_, err := r.db.Exec("DELETE FROM stores WHERE id = $1", id)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func (r *storeRepositoryImpl) FindAll() ([]entities.Store, error) {
	rows, err := r.db.Query(`SELECT 
    id, 
    name,
    description,
    created_at, 
    updated_at 
	FROM stores`)

	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	defer rows.Close()

	var stores []entities.Store
	for rows.Next() {
		var store entities.Store
		if err := rows.Scan(&store.ID,
			&store.Name,
			&store.Description,
			&store.CreatedAt,
			&store.UpdatedAt); err != nil {
			return nil, tracerr.Wrap(err)
		}
		stores = append(stores, store)
	}

	if err = rows.Err(); err != nil {
		return nil, tracerr.Wrap(err)
	}

	return stores, nil
}

func (r *storeRepositoryImpl) FindStoresByParams(params entities.StoreSearchParams) ([]entities.Store, *uint64, error) {
	query, args := r.buildQuery(params)
	stores, lastCursor, err := r.executeAndProcessQuery(query, args)
	if err != nil {
		return nil, nil, tracerr.Wrap(fmt.Errorf("error executing or processing query: %w", err))
	}

	return stores, lastCursor, nil
}

// buildQuery генерирует SQL-запрос и параметры.
func (r *storeRepositoryImpl) buildQuery(params entities.StoreSearchParams) (string, []interface{}) {
	var query string
	var conditions []string
	var args []interface{}

	// Если указан CategoryID, добавляем JOIN
	if params.CategoryID != nil {
		query = `SELECT s.id, s.name, s.description, s.created_at, s.updated_at 
			FROM stores s
			INNER JOIN store_categories sc ON s.id = sc.store_id`
		conditions = append(conditions, fmt.Sprintf("sc.category_id = $%d", len(args)+1))
		args = append(args, *params.CategoryID)
	} else {
		query = `SELECT id, name, description, created_at, updated_at 
			FROM stores`
	}

	if params.Name != nil {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Name+"%")
	}

	if params.Cursor != nil {
		conditions = append(conditions, fmt.Sprintf("id > $%d", len(args)+1))
		args = append(args, *params.Cursor)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY id ASC"
	if params.Limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, *params.Limit)
	}

	return query, args
}

// executeAndProcessQuery выполняет запрос и обрабатывает результаты.
func (r *storeRepositoryImpl) executeAndProcessQuery(
	query string,
	args []interface{},
) ([]entities.Store, *uint64, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, nil, tracerr.Wrap(fmt.Errorf("error executing query: %w", err))
	}
	defer rows.Close()

	var stores []entities.Store
	var lastCursor *uint64

	for rows.Next() {
		var store entities.Store
		if err = rows.Scan(
			&store.ID,
			&store.Name,
			&store.Description,
			&store.CreatedAt,
			&store.UpdatedAt,
		); err != nil {
			return nil, nil, tracerr.Wrap(fmt.Errorf("error scanning row: %w", err))
		}
		stores = append(stores, store)
		lastCursor = &store.ID
	}

	if err = rows.Err(); err != nil {
		return nil, nil, tracerr.Wrap(fmt.Errorf("error during rows iteration: %w", err))
	}

	return stores, lastCursor, nil
}

func (r *storeRepositoryImpl) AttachCategory(storeID, categoryID uint64) error {
	_, err := r.db.Exec(`INSERT INTO store_categories (store_id, category_id) VALUES ($1, $2)`, storeID, categoryID)
	if err != nil {
		return tracerr.Wrap(fmt.Errorf("failed to add category to store: %w", err))
	}
	return nil
}

func (r *storeRepositoryImpl) IsCategoryAttached(storeID, categoryID uint64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM store_categories WHERE store_id = $1 AND category_id = $2)`

	err := r.db.QueryRow(query, storeID, categoryID).Scan(&exists)
	if err != nil {
		return false, tracerr.Wrap(err)
	}

	return exists, nil
}

func (r *storeRepositoryImpl) DetachCategory(storeID, categoryID uint64) error {
	_, err := r.db.Exec(`DELETE FROM store_categories WHERE store_id = $1 AND category_id = $2`, storeID, categoryID)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func (r *storeRepositoryImpl) IsUserStoreAdmin(storeID, uid uint64) (bool, error) {
	var isAdmin bool
	query := `SELECT EXISTS(SELECT 1 FROM store_roles WHERE store_id = $1 AND user_id = $2)`

	err := r.db.QueryRow(query, storeID, uid).Scan(&isAdmin)
	if err != nil {
		return false, tracerr.Wrap(fmt.Errorf("failed to check if user is admin: %w", err))
	}

	return isAdmin, nil
}

func (r *storeRepositoryImpl) AddUserStoreAdmin(storeID, uid uint64, owner bool) error {
	query := `INSERT INTO store_roles (
                         store_id, 
                         user_id, 
                         is_owner) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`

	_, err := r.db.Exec(query, storeID, uid, owner)
	if err != nil {
		return tracerr.Wrap(fmt.Errorf("failed to add user as admin to store: %w", err))
	}

	return nil
}
