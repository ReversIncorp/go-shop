package repository

import (
	"errors"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
	errorHandling "marketplace/pkg/error_handling"
	"time"

	"github.com/ztrue/tracerr"
	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepositoryImpl{
		db: db,
	}
}

func (r *productRepositoryImpl) Save(product *entities.Product) error {
	product.CreatedAt = time.Now()
	product.UpdatedAt = product.CreatedAt
	if err := r.db.Create(product).Error; err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func (r *productRepositoryImpl) FindByID(id uint64) (entities.Product, error) {
	var product entities.Product
	err := r.db.Preload("Category").
	Preload("Store").First(&product, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.Product{}, errorHandling.ErrProductNotFound
	}
	if err != nil {
		return entities.Product{}, tracerr.Wrap(err)
	}
	return product, nil
}

func (r *productRepositoryImpl) Update(product entities.Product) error {
	var existing entities.Product
	err := r.db.First(&existing, product.ID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errorHandling.ErrProductNotFound
	}
	if err != nil {
		return tracerr.Wrap(err)
	}
	product.UpdatedAt = time.Now()
	if err := r.db.Model(&product).Updates(map[string]interface{}{
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"quantity":    product.Quantity,
		"category_id": product.CategoryID,
		"store_id":    product.StoreID,
		"updated_at":  product.UpdatedAt,
	}).Error; err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func (r *productRepositoryImpl) Delete(id uint64) error {
	if err := r.db.Delete(&entities.Product{}, id).Error; err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func (r *productRepositoryImpl) FindProductsByParams(params entities.ProductSearchParams) ([]entities.Product, *uint64, error) {
	db := r.db.Model(&entities.Product{})
	if params.StoreID != nil {
		db = db.Where("store_id = ?", *params.StoreID)
	}
	if params.CategoryID != nil {
		db = db.Where("category_id = ?", *params.CategoryID)
	}
	if params.MinPrice != nil {
		db = db.Where("price >= ?", *params.MinPrice)
	}
	if params.MaxPrice != nil {
		db = db.Where("price <= ?", *params.MaxPrice)
	}
	if params.Name != nil {
		db = db.Where("name ILIKE ?", "%"+*params.Name+"%")
	}
	if params.Cursor != nil {
		db = db.Where("id > ?", *params.Cursor)
	}
	if params.Limit != nil {
		db = db.Limit(int(*params.Limit))
	}
	var products []entities.Product
	if err := db.
		Preload("Category").
		Preload("Store").
		Order("id ASC").
		Find(&products).Error; err != nil {
		return nil, nil, tracerr.Wrap(err)
	}
	var lastCursor *uint64
	if len(products) > 0 {
		lastCursor = &products[len(products)-1].ID
	}
	return products, lastCursor, nil
}

func (r *productRepositoryImpl) IsProductBelongsToStore(productID, storeID uint64) (bool, error) {
	var count int64
	if err := r.db.Model(&entities.Product{}).Where("id = ? AND store_id = ?", productID, storeID).Count(&count).Error; err != nil {
		return false, tracerr.Wrap(err)
	}
	return count > 0, nil
}
