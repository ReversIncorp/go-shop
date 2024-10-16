package repository

import (
	"errors"
	"gorm.io/gorm"
	"marketplace/internal/domain/entities"
	repository2 "marketplace/internal/domain/repository"
	"sync"
	"time"
)

type productRepositoryImpl struct {
	db *gorm.DB
	mu sync.Mutex
}

func NewProductRepository(db *gorm.DB) repository2.ProductRepository {
	return &productRepositoryImpl{
		db: db,
	}
}

func (r *productRepositoryImpl) Save(product entities.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var existingProduct entities.Product
	if err := r.db.Where("id = ?", product.ID).First(&existingProduct).Error; err == nil {
		return errors.New("product already exists")
	}

	product.CreatedAt = time.Now()
	product.UpdatedAt = product.CreatedAt

	if err := r.db.Create(&product).Error; err != nil {
		return err
	}

	return nil
}

func (r *productRepositoryImpl) FindByID(id uint64) (entities.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var product entities.Product
	if err := r.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Product{}, errors.New("product not found")
		}
		return entities.Product{}, err
	}

	return product, nil
}

func (r *productRepositoryImpl) Update(product entities.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var existingProduct entities.Product
	if err := r.db.First(&existingProduct, product.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	product.UpdatedAt = time.Now()

	if err := r.db.Save(&product).Error; err != nil {
		return err
	}

	return nil
}

func (r *productRepositoryImpl) Delete(id uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.db.Delete(&entities.Product{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	return nil
}

func (r *productRepositoryImpl) FindAllByStore(storeID uint64) ([]entities.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var products []entities.Product
	if err := r.db.Where("store_id = ?", storeID).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}
