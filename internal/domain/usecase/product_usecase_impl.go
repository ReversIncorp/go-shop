package usecase

import (
	"marketplace/internal/domain/models"
	"marketplace/internal/domain/repository"
)

// productUseCaseImpl реализует интерфейс ProductUseCase
type productUseCaseImpl struct {
	productRepo repository.ProductRepository
}

// NewProductUseCase создает новый экземпляр ProductUseCase
func NewProductUseCase(productRepo repository.ProductRepository) ProductUseCase {
	return &productUseCaseImpl{productRepo: productRepo}
}

// CreateProduct создает новый продукт
func (p *productUseCaseImpl) CreateProduct(product models.Product) error {
	return p.productRepo.Save(product)
}

// GetProductByID получает продукт по ID
func (p *productUseCaseImpl) GetProductByID(id string) (models.Product, error) {
	return p.productRepo.FindByID(id)
}

// UpdateProduct обновляет существующий продукт
func (p *productUseCaseImpl) UpdateProduct(product models.Product) error {
	return p.productRepo.Update(product)
}

// DeleteProduct удаляет продукт по ID
func (p *productUseCaseImpl) DeleteProduct(id string) error {
	return p.productRepo.Delete(id)
}

// GetProductsByStore получает все продукты по ID магазина
func (p *productUseCaseImpl) GetProductsByStore(storeID string) ([]models.Product, error) {
	return p.productRepo.FindAllByStore(storeID)
}
