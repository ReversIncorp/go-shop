package repository

import (
	"errors"
	"marketplace/internal/domain/entities"
	repository2 "marketplace/internal/domain/repository"
	"sync"
	"time"
)

type inMemoryProductRepository struct {
	products map[uint64]entities.Product
	mu       sync.Mutex
}

func NewProductRepository() repository2.ProductRepository {
	return &inMemoryProductRepository{
		products: make(map[uint64]entities.Product),
	}
}

func (r *inMemoryProductRepository) Save(product entities.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Если продукт уже существует, вернем ошибку
	if _, exists := r.products[product.ID]; exists {
		return errors.New("product already exists")
	}

	// Устанавливаем время создания и обновления
	product.CreatedAt = time.Now()
	product.UpdatedAt = product.CreatedAt

	r.products[product.ID] = product
	return nil
}

func (r *inMemoryProductRepository) FindByID(id uint64) (entities.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	product, exists := r.products[id]
	if !exists {
		return entities.Product{}, errors.New("product not found")
	}

	return product, nil
}

func (r *inMemoryProductRepository) Update(product entities.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.products[product.ID]
	if !exists {
		return errors.New("product not found")
	}

	// Обновляем время изменения
	product.UpdatedAt = time.Now()

	r.products[product.ID] = product
	return nil
}

func (r *inMemoryProductRepository) Delete(id uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.products[id]
	if !exists {
		return errors.New("product not found")
	}

	delete(r.products, id)
	return nil
}

func (r *inMemoryProductRepository) FindAllByStore(storeID uint64) ([]entities.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var products []entities.Product
	for _, product := range r.products {
		if product.StoreID == storeID {
			products = append(products, product)
		}
	}

	return products, nil
}
