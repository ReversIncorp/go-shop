package productUsecase

import (
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
	errorResponses "marketplace/pkg/errors"
)

// ProductUseCase реализует интерфейс ProductUseCase
type ProductUseCase struct {
	productRepo repository.ProductRepository
	storeRepo   repository.StoreRepository
}

// NewProductUseCase создает новый экземпляр ProductUseCase
func NewProductUseCase(
	productRepo repository.ProductRepository,
	storeRepo repository.StoreRepository,
) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
		storeRepo:   storeRepo,
	}
}

// CreateProduct создает новый продукт
func (p *ProductUseCase) CreateProduct(product entities.Product) error {
	categoryBelongs, err := p.storeRepo.IsCategoryAttached(product.StoreID, product.CategoryID)
	if err != nil || !categoryBelongs {
		return errorResponses.ErrCategoryNotAttached
	}

	err = p.productRepo.Save(product)
	if err != nil {
		return errorResponses.ErrInternalServerError
	}

	return nil
}

// GetProductByID получает продукт по ID
func (p *ProductUseCase) GetProductByID(id uint64) (entities.Product, error) {
	return p.productRepo.FindByID(id)
}

// UpdateProduct обновляет существующий продукт
func (p *ProductUseCase) UpdateProduct(product entities.Product) error {
	categoryBelongs, err := p.storeRepo.IsCategoryAttached(product.StoreID, product.CategoryID)
	if err != nil || !categoryBelongs {
		return errorResponses.ErrCategoryNotAttached
	}

	productBelongs, err := p.productRepo.IsProductBelongsToStore(product.ID, product.StoreID)
	if err != nil || !productBelongs {
		return errorResponses.ErrProductNotBelongsToStore
	}

	err = p.productRepo.Save(product)
	if err != nil {
		return errorResponses.ErrInternalServerError
	}

	return nil
}

// DeleteProduct удаляет продукт по ID
func (p *ProductUseCase) DeleteProduct(id uint64) error {
	product, err := p.productRepo.FindByID(id)
	if err != nil {
		return errorResponses.ErrProductNotFound
	}

	productBelongs, err := p.productRepo.IsProductBelongsToStore(product.ID, product.StoreID)
	if err != nil || !productBelongs {
		return errorResponses.ErrProductNotBelongsToStore
	}

	err = p.productRepo.Delete(id)
	if err != nil {
		return errorResponses.ErrInternalServerError
	}

	return nil
}

// GetProductsByFilters получает все продукты по фильтрам
func (p *ProductUseCase) GetProductsByFilters(filters entities.ProductSearchParams) ([]entities.Product, error) {
	return p.productRepo.FindProductsByParams(filters)
}
