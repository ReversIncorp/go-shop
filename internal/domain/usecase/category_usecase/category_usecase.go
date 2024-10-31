package categoryUsecase

import (
	"errors"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
)

type CategoryUseCase struct {
	categoryRepo repository.CategoryRepository
	userRepo     repository.UserRepository
}

func NewCategoryUseCase(
	categoryRepo repository.CategoryRepository,
	userRepo repository.UserRepository,
) *CategoryUseCase {
	return &CategoryUseCase{
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
	}
}

func (c *CategoryUseCase) CreateCategory(category entities.Category, uid uint64) error {
	userData, err := c.userRepo.FindByID(uid)
	if err != nil || !userData.IsSeller {
		return errors.New("user is not seller")
	}

	return c.categoryRepo.Save(category)
}

func (c *CategoryUseCase) UpdateCategory(category entities.Category, uid uint64) error {
	userData, err := c.userRepo.FindByID(uid)
	if err != nil || !userData.IsSeller {
		return errors.New("user is not seller")
	}

	categoryExists, err := c.categoryRepo.IsExist(category.ID)
	if err != nil || !categoryExists {
		return errors.New("category not found")
	}

	return c.categoryRepo.Update(category)
}

func (c *CategoryUseCase) DeleteCategory(id, uid uint64) error {
	userData, err := c.userRepo.FindByID(uid)
	if err != nil || !userData.IsSeller {
		return errors.New("user is not seller")
	}

	return c.categoryRepo.Delete(id)
}

func (c *CategoryUseCase) GetAllCategoriesByStore(id uint64) ([]entities.Category, error) {
	return c.categoryRepo.FindAllByStore(id)
}
