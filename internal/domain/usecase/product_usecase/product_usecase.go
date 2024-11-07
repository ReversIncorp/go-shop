package productUsecase

import (
	"errors"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
)

// ProductUseCase реализует интерфейс ProductUseCase
type ProductUseCase struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

// NewProductUseCase создает новый экземпляр ProductUseCase
func NewProductUseCase(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
) *ProductUseCase {
	return &ProductUseCase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateProduct создает новый продукт
func (p *ProductUseCase) CreateProduct(product entities.Product) error {
	categoryBelongs, err := p.categoryRepo.IsBelongsToStore(product.CategoryID, product.StoreID)
	if err != nil || !categoryBelongs {
		return errors.New("category not found or not belongs this store")
	}

	return p.productRepo.Save(product)
}

// GetProductByID получает продукт по ID
func (p *ProductUseCase) GetProductByID(id uint64) (entities.Product, error) {
	return p.productRepo.FindByID(id)
}

// UpdateProduct обновляет существующий продукт
func (p *ProductUseCase) UpdateProduct(product entities.Product) error {
	categoryBelongs, err := p.categoryRepo.IsBelongsToStore(product.CategoryID, product.StoreID)
	if err != nil || !categoryBelongs {
		return errors.New("category not found or not belongs this store")
	}

	productBelongs, err := p.productRepo.IsProductBelongsToStore(product.ID, product.StoreID)
	if err != nil || !productBelongs {
		return errors.New("product not found or not belongs this store")
	}

	return p.productRepo.Update(product)
}

// DeleteProduct удаляет продукт по ID
func (p *ProductUseCase) DeleteProduct(id uint64) error {
	product, err := p.productRepo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	productBelongs, err := p.productRepo.IsProductBelongsToStore(product.ID, product.StoreID)
	if err != nil || !productBelongs {
		return errors.New("product not belongs this store")
	}

	return p.productRepo.Delete(id)
}

// GetProductsByFilters получает все продукты по фильтрам
func (p *ProductUseCase) GetProductsByFilters(filters entities.ProductSearchParams) ([]entities.Product, error) {
	return p.productRepo.FindProductsByParams(filters)
}
