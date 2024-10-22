package usecase

import (
	"errors"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
)

// ProductUseCase реализует интерфейс ProductUseCase
type ProductUseCase struct {
	productRepo  repository.ProductRepository
	storeRepo    repository.StoreRepository
	userRepo     repository.UserRepository
	categoryRepo repository.CategoryRepository
}

// NewProductUseCase создает новый экземпляр ProductUseCase
func NewProductUseCase(productRepo repository.ProductRepository, storeRepo repository.StoreRepository, userRepo repository.UserRepository, categoryRepo repository.CategoryRepository) *ProductUseCase {
	return &ProductUseCase{productRepo: productRepo, storeRepo: storeRepo, userRepo: userRepo, categoryRepo: categoryRepo}
}

// CreateProduct создает новый продукт
func (p *ProductUseCase) CreateProduct(product entities.Product, uid int64) error {
	storeExists, err := p.storeRepo.IsExist(product.StoreID)
	if err != nil || !storeExists {
		return errors.New("store not found")
	}

	categoryBelongs, err := p.categoryRepo.IsBelongsToStore(product.CategoryID, product.StoreID)
	if err != nil || !categoryBelongs {
		return errors.New("category not found or not belongs this store")
	}

	isOwner, err := p.userRepo.IsOwnsStore(uid, product.StoreID)
	if err != nil || !isOwner {
		return errors.New("user does not owning this store")
	}

	return p.productRepo.Save(product)
}

// GetProductByID получает продукт по ID
func (p *ProductUseCase) GetProductByID(id int64) (entities.Product, error) {
	return p.productRepo.FindByID(id)
}

// UpdateProduct обновляет существующий продукт
func (p *ProductUseCase) UpdateProduct(product entities.Product, uid int64) error {
	storeExists, err := p.storeRepo.IsExist(product.StoreID)
	if err != nil || !storeExists {
		return errors.New("store not found")
	}

	categoryBelongs, err := p.categoryRepo.IsBelongsToStore(product.CategoryID, product.StoreID)
	if err != nil || !categoryBelongs {
		return errors.New("category not found or not belongs this store")
	}

	isOwner, err := p.userRepo.IsOwnsStore(uid, product.StoreID)
	if err != nil || !isOwner {
		return errors.New("user does not owning this store")
	}

	return p.productRepo.Update(product)
}

// DeleteProduct удаляет продукт по ID
func (p *ProductUseCase) DeleteProduct(id int64, uid int64) error {
	product, err := p.productRepo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	isOwner, err := p.userRepo.IsOwnsStore(uid, product.StoreID)
	if err != nil || !isOwner {
		return errors.New("user does not owning this store")
	}

	return p.productRepo.Delete(id)
}

// GetProductsByStore получает все продукты по ID магазина
func (p *ProductUseCase) GetProductsByStore(storeID int64) ([]entities.Product, error) {
	return p.productRepo.FindAllByStore(storeID)
}

// GetProductsByStoreAndCategory получает все продукты по ID магазина
func (p *ProductUseCase) GetProductsByStoreAndCategory(storeID int64, categoryID int64) ([]entities.Product, error) {
	return p.productRepo.FindAllByStoreAndCategory(storeID, categoryID)
}
