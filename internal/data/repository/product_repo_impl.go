package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
	"strings"
	"time"
)

type productRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) repository.ProductRepository {
	return &productRepositoryImpl{
		db: db,
	}
}

func (r *productRepositoryImpl) Save(product entities.Product) error {
	product.CreatedAt = time.Now()
	product.UpdatedAt = product.CreatedAt

	_, err := r.db.Exec(`INSERT INTO products 
    (name, 
     description, 
     price, 
     quantity, 
     category_id, 
     store_id, 
     created_at, 
     updated_at) 
                        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		product.Name,
		product.Description,
		product.Price,
		product.Quantity,
		product.CategoryID,
		product.StoreID,
		product.CreatedAt,
		product.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *productRepositoryImpl) FindByID(id uint64) (entities.Product, error) {
	var product entities.Product
	err := r.db.QueryRow(`SELECT 
    id, 
    name, 
    description, 
    price, 
    quantity, 
    category_id, 
    store_id, 
    created_at, 
    updated_at 
	                      FROM products WHERE id = $1`, id).
		Scan(&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Quantity,
			&product.CategoryID,
			&product.StoreID,
			&product.CreatedAt,
			&product.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return entities.Product{}, errors.New("product not found")
	}
	if err != nil {
		return entities.Product{}, err
	}

	return product, nil
}

func (r *productRepositoryImpl) Update(product entities.Product) error {
	var existingProductID int64
	err := r.db.QueryRow("SELECT id FROM products WHERE id = $1", product.ID).Scan(&existingProductID)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("product not found")
	}
	if err != nil {
		return err
	}

	product.UpdatedAt = time.Now()

	_, err = r.db.Exec(`UPDATE products 
                        SET name = $1, 
                            description = $2, 
                            price = $3, 
                            quantity = $4, 
                            category_id = $5, 
                            store_id = $6, 
                            updated_at = $7 
                        WHERE id = $8`,
		product.Name,
		product.Description,
		product.Price,
		product.Quantity,
		product.CategoryID,
		product.StoreID,
		product.UpdatedAt,
		product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepositoryImpl) Delete(id uint64) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepositoryImpl) FindProductsByParams(params entities.ProductSearchParams) ([]entities.Product, error) {
	query := `SELECT 
    id, 
    name,
    description,
    price, 
    quantity, 
    category_id,
    store_id,
    created_at, 
    updated_at 
	FROM products`

	var conditions []string
	var args []interface{}

	if params.StoreID != nil {
		conditions = append(conditions, fmt.Sprintf("store_id = $%d", len(args)+1))
		args = append(args, *params.StoreID)
	}

	if params.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", len(args)+1))
		args = append(args, *params.CategoryID)
	}

	if params.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price >= $%d", len(args)+1))
		args = append(args, *params.MinPrice)
	}

	if params.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price <= $%d", len(args)+1))
		args = append(args, *params.MaxPrice)
	}

	if params.Name != nil {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Name+"%")
	}

	// Добавляем условия, если они есть
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var products []entities.Product
	for rows.Next() {
		var product entities.Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Quantity,
			&product.CategoryID,
			&product.StoreID,
			&product.CreatedAt,
			&product.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return products, nil
}
