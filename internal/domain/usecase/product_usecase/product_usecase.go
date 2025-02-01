package productusecase

import (
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
	errorHandling "marketplace/pkg/error_handling"
	"marketplace/pkg/utils"

	"github.com/ztrue/tracerr"
)

// ProductUseCase реализует интерфейс ProductUseCase.
type ProductUseCase struct {
	productRepo repository.ProductRepository
	storeRepo   repository.StoreRepository
}

// NewProductUseCase создает новый экземпляр ProductUseCase.
func NewProductUseCase(
	productRepo repository.ProductRepository,
	storeRepo repository.StoreRepository,
) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
		storeRepo:   storeRepo,
	}
}

// CreateProduct создает новый продукт.
func (p *ProductUseCase) CreateProduct(product entities.Product) error {
	categoryBelongs, err := p.storeRepo.IsCategoryAttached(product.StoreID, product.CategoryID)
	if err != nil || !categoryBelongs {
		return errorHandling.ErrCategoryNotAttached
	}

	err = p.productRepo.Save(product)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

// GetProductByID получает продукт по ID.
func (p *ProductUseCase) GetProductByID(id uint64) (entities.Product, error) {
	product, err := p.productRepo.FindByID(id)
	if err != nil {
		return product, utils.GetHttpErrorOrTracerrError(err)
	}
	return product, nil
}

// UpdateProduct обновляет существующий продукт.
func (p *ProductUseCase) UpdateProduct(product entities.Product) error {
	categoryBelongs, err := p.storeRepo.IsCategoryAttached(product.StoreID, product.CategoryID)
	if err != nil || !categoryBelongs {
		return errorHandling.ErrCategoryNotAttached
	}

	productBelongs, err := p.productRepo.IsProductBelongsToStore(product.ID, product.StoreID)
	if err != nil || !productBelongs {
		return errorHandling.ErrProductNotBelongsToStore
	}

	err = p.productRepo.Save(product)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

// DeleteProduct удаляет продукт по ID.
func (p *ProductUseCase) DeleteProduct(id uint64) error {
	product, err := p.productRepo.FindByID(id)
	if err != nil {
		return errorHandling.ErrProductNotFound
	}

	productBelongs, err := p.productRepo.IsProductBelongsToStore(product.ID, product.StoreID)
	if err != nil || !productBelongs {
		return errorHandling.ErrProductNotBelongsToStore
	}

	err = p.productRepo.Delete(id)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

// GetProductsByFilters получает все продукты по фильтрам.
func (p *ProductUseCase) GetProductsByFilters(
	filters entities.ProductSearchParams,
) ([]entities.Product, *uint64, error) {
	products, cursor, err := p.productRepo.FindProductsByParams(filters)
	if err != nil {
		return nil, nil, tracerr.Wrap(err)
	}
	return products, cursor, err
}
