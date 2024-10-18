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
func (p *ProductUseCase) CreateProduct(product entities.Product, uid int64) error {
	return p.productRepo.Save(product, uid)
}

// GetProductByID получает продукт по ID
func (p *ProductUseCase) GetProductByID(id int64) (entities.Product, error) {
	return p.productRepo.FindByID(id)
}

// UpdateProduct обновляет существующий продукт
func (p *ProductUseCase) UpdateProduct(product entities.Product, uid int64) error {
	return p.productRepo.Update(product, uid)
}

// DeleteProduct удаляет продукт по ID
func (p *ProductUseCase) DeleteProduct(id int64, uid int64) error {
	return p.productRepo.Delete(id, uid)
}

// GetProductsByStore получает все продукты по ID магазина
func (p *ProductUseCase) GetProductsByStore(storeID int64) ([]entities.Product, error) {
	return p.productRepo.FindAllByStore(storeID)
}

// GetProductsByStoreAndCategory получает все продукты по ID магазина
func (p *ProductUseCase) GetProductsByStoreAndCategory(storeID int64, categoryID int64) ([]entities.Product, error) {
	return p.productRepo.FindAllByStoreAndCategory(storeID, categoryID)
}
