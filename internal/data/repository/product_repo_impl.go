package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
	errorHandling "marketplace/pkg/error_handling"
	"strings"
	"time"

	"github.com/ztrue/tracerr"
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
		return tracerr.Wrap(err)
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
		return entities.Product{}, errorHandling.ErrProductNotFound
	}
	if err != nil {
		return entities.Product{}, tracerr.Wrap(err)
	}

	return product, nil
}

func (r *productRepositoryImpl) Update(product entities.Product) error {
	var existingProductID int64
	err := r.db.QueryRow("SELECT id FROM products WHERE id = $1", product.ID).Scan(&existingProductID)
	if errors.Is(err, sql.ErrNoRows) {
		return errorHandling.ErrProductNotFound
	}
	if err != nil {
		return tracerr.Wrap(err)
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
		return tracerr.Wrap(err)
	}

	return nil
}

func (r *productRepositoryImpl) Delete(id uint64) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func (r *productRepositoryImpl) FindProductsByParams(
	params entities.ProductSearchParams,
) ([]entities.Product, *uint64, error) {
	query, args := r.buildProductQuery(params)
	rows, err := r.executeQuery(query, args)
	if err != nil {
		return nil, nil, tracerr.Wrap(err)
	}
	defer rows.Close()

	products, lastCursor, err := r.processProductRows(rows, params.Limit)
	if err != nil {
		return nil, nil, tracerr.Wrap(err)
	}

	return products, lastCursor, nil
}

// Генерирует SQL-запрос и параметры.
func (r *productRepositoryImpl) buildProductQuery(params entities.ProductSearchParams) (string, []interface{}) {
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

	// Фильтры по параметрам
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

	// Курсор для пагинации
	if params.Cursor != nil {
		conditions = append(conditions, fmt.Sprintf("id > $%d", len(args)+1))
		args = append(args, *params.Cursor)
	}

	// Добавляем условия в запрос
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Сортировка по ID
	query += " ORDER BY id ASC"

	// Лимит
	if params.Limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, *params.Limit)
	}

	return query, args
}

// Выполняет SQL-запрос.
func (r *productRepositoryImpl) executeQuery(query string, args []interface{}) (*sql.Rows, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return rows, nil
}

// Обрабатывает результаты запроса.
func (r *productRepositoryImpl) processProductRows(rows *sql.Rows, limit *uint64) ([]entities.Product, *uint64, error) {
	var products []entities.Product
	var lastCursor uint64

	// Читаем строки
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
			return nil, nil, tracerr.Wrap(err)
		}
		products = append(products, product)
		lastCursor = product.ID
	}

	// Проверяем ошибки чтения строк
	if err := rows.Err(); err != nil {
		return nil, nil, tracerr.Wrap(err)
	}

	// Если данных меньше лимита, следующего курсора нет
	if limit != nil && uint64(len(products)) < *limit {
		return products, nil, nil
	}

	return products, &lastCursor, nil
}

func (r *productRepositoryImpl) IsProductBelongsToStore(productID, storeID uint64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(
		SELECT 1 
		FROM products 
		WHERE id = $1 AND store_id = $2
	)`

	err := r.db.QueryRow(query, productID, storeID).Scan(&exists)
	if err != nil {
		return false, tracerr.Wrap(err)
	}

	return exists, nil
}
