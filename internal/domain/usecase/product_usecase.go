package usecase

import (
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
)

// ProductUseCase реализует интерфейс ProductUseCase
type ProductUseCase struct {
	productRepo repository.ProductRepository
}

// NewProductUseCase создает новый экземпляр ProductUseCase
func NewProductUseCase(productRepo repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{productRepo: productRepo}
}

// CreateProduct создает новый продукт
func (p *ProductUseCase) CreateProduct(product entities.Product) error {
	return p.productRepo.Save(product)
}

// GetProductByID получает продукт по ID
func (p *ProductUseCase) GetProductByID(id string) (entities.Product, error) {
	return p.productRepo.FindByID(id)
}

// UpdateProduct обновляет существующий продукт
func (p *ProductUseCase) UpdateProduct(product entities.Product) error {
	return p.productRepo.Update(product)
}

// DeleteProduct удаляет продукт по ID
func (p *ProductUseCase) DeleteProduct(id string) error {
	return p.productRepo.Delete(id)
}

// GetProductsByStore получает все продукты по ID магазина
func (p *ProductUseCase) GetProductsByStore(storeID string) ([]entities.Product, error) {
	return p.productRepo.FindAllByStore(storeID)
}
