package usecase

import (
	"errors"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
)

type CategoryUseCase struct {
	categoryRepo repository.CategoryRepository
	storeRepo    repository.StoreRepository
	userRepo     repository.UserRepository
}

func NewCategoryUseCase(categoryRepo repository.CategoryRepository, storeRepo repository.StoreRepository, userRepo repository.UserRepository) *CategoryUseCase {
	return &CategoryUseCase{categoryRepo: categoryRepo, storeRepo: storeRepo, userRepo: userRepo}
}

func (c *CategoryUseCase) CreateCategory(category entities.Category, uid int64) error {
	storeExists, err := c.storeRepo.IsExist(category.StoreID)
	if err != nil || !storeExists {
		return errors.New("store not found")
	}

	isOwner, err := c.userRepo.IsOwnsStore(uid, category.StoreID)
	if err != nil || !isOwner {
		return errors.New("user does not owning this store")
	}

	return c.categoryRepo.Save(category)
}

func (c *CategoryUseCase) GetCategoryByID(id int64) (entities.Category, error) {
	return c.categoryRepo.FindByID(id)
}

func (c *CategoryUseCase) UpdateCategory(category entities.Category, uid int64) error {
	categoryExists, err := c.categoryRepo.IsExist(category.ID)
	if err != nil || !categoryExists {
		return errors.New("category not found")
	}

	storeExists, err := c.storeRepo.IsExist(category.StoreID)
	if err != nil || !storeExists {
		return errors.New("store not found")
	}

	isOwner, err := c.userRepo.IsOwnsStore(uid, category.StoreID)
	if err != nil || !isOwner {
		return errors.New("user does not owning this store")
	}

	return c.categoryRepo.Update(category)
}

func (c *CategoryUseCase) DeleteCategory(id int64, uid int64) error {
	category, err := c.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	isOwner, err := c.userRepo.IsOwnsStore(uid, category.StoreID)
	if err != nil || !isOwner {
		return errors.New("user does not owning this store")
	}

	return c.categoryRepo.Delete(id)
}

func (c *CategoryUseCase) GetAllCategoriesByStore(id int64) ([]entities.Category, error) {
	return c.categoryRepo.FindAllByStore(id)
}
