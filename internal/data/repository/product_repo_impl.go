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

type productRepositoryImpl struct {
	db *sql.DB
	mu sync.Mutex
}

func NewProductRepository(db *sql.DB) repository2.ProductRepository {
	return &productRepositoryImpl{
		db: db,
	}
}

func (r *productRepositoryImpl) Save(product entities.Product, uid int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storeExists, err := utils.CheckStoreExists(r.db, product.StoreID)
	if err != nil || !storeExists {
		return errors.New("store not found")
	}

	isOwner, err := utils.CheckUserOwnsStore(r.db, uid, product.StoreID)
	if err != nil || !isOwner {
		return errors.New("user does not own this store")
	}

	belongs, err := utils.CheckCategoryBelongsToStore(r.db, product.CategoryID, product.StoreID)
	if err != nil || !belongs {
		return errors.New("category does not belong to this store")
	}

	product.CreatedAt = time.Now()
	product.UpdatedAt = product.CreatedAt

	_, err = r.db.Exec(`INSERT INTO products (name, description, price, quantity, category_id, store_id, created_at, updated_at) 
                        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		product.Name, product.Description, product.Price, product.Quantity, product.CategoryID, product.StoreID, product.CreatedAt, product.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepositoryImpl) FindByID(id int64) (entities.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var product entities.Product
	err := r.db.QueryRow(`SELECT id, name, description, price, quantity, category_id, store_id, created_at, updated_at 
	                      FROM products WHERE id = $1`, id).
		Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Quantity, &product.CategoryID, &product.StoreID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Product{}, errors.New("product not found")
		}
		return entities.Product{}, err
	}

	return product, nil
}

func (r *productRepositoryImpl) Update(product entities.Product, uid int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var existingProductID int64
	err := r.db.QueryRow("SELECT id FROM products WHERE id = $1", product.ID).Scan(&existingProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("product not found")
		}
		return err
	}

	storeExists, err := utils.CheckStoreExists(r.db, product.StoreID)
	if err != nil || !storeExists {
		return errors.New("store not found")
	}

	isOwner, err := utils.CheckUserOwnsStore(r.db, uid, product.StoreID)
	if err != nil || !isOwner {
		return errors.New("user does not own this store")
	}

	belongs, err := utils.CheckCategoryBelongsToStore(r.db, product.CategoryID, product.StoreID)
	if err != nil || !belongs {
		return errors.New("category does not belong to this store")
	}

	product.UpdatedAt = time.Now()

	_, err = r.db.Exec(`UPDATE products 
                        SET name = $1, description = $2, price = $3, quantity = $4, category_id = $5, store_id = $6, updated_at = $7 
                        WHERE id = $8`,
		product.Name, product.Description, product.Price, product.Quantity, product.CategoryID, product.StoreID, product.UpdatedAt, product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepositoryImpl) Delete(id int64, uid int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var existingProductID int64
	var existingProductStoreID int64
	err := r.db.QueryRow(`SELECT id, store_id FROM products WHERE id = $1`, id).
		Scan(&existingProductID, &existingProductStoreID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("product not found")
		}
		return err
	}

	storeExists, err := utils.CheckStoreExists(r.db, existingProductStoreID)
	if err != nil || !storeExists {
		return errors.New("store not found")
	}

	isOwner, err := utils.CheckUserOwnsStore(r.db, uid, existingProductStoreID)
	if err != nil || !isOwner {
		return errors.New("user does not own this store")
	}

	_, err = r.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepositoryImpl) FindAllByStore(storeID int64) ([]entities.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.Query(`SELECT id, name, description, price, quantity, category_id, store_id, created_at, updated_at 
	                         FROM products WHERE store_id = $1`, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []entities.Product
	for rows.Next() {
		var product entities.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Quantity, &product.CategoryID, &product.StoreID, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepositoryImpl) FindAllByStoreAndCategory(storeID int64, categoryID int64) ([]entities.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	belongs, err := utils.CheckCategoryBelongsToStore(r.db, categoryID, storeID)
	if err != nil || !belongs {
		return nil, errors.New("category does not belong to this store")
	}

	rows, err := r.db.Query(`SELECT id, name, description, price, quantity, category_id, store_id, created_at, updated_at 
	                         FROM products WHERE store_id = $1 AND category_id = $2`, storeID, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []entities.Product
	for rows.Next() {
		var product entities.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Quantity, &product.CategoryID, &product.StoreID, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
