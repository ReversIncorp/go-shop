package usecase

import (
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
)

type CategoryUseCase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUseCase(categoryRepo repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{categoryRepo: categoryRepo}
}

func (c *CategoryUseCase) CreateCategory(category entities.Category) error {
	return c.categoryRepo.Save(category)
}

func (c *CategoryUseCase) GetCategoryByID(id uint64) (entities.Category, error) {
	return c.categoryRepo.FindByID(id)
}

func (c *CategoryUseCase) UpdateCategory(category entities.Category) error {
	return c.categoryRepo.Update(category)
}

func (c *CategoryUseCase) DeleteCategory(id uint64) error {
	return c.categoryRepo.Delete(id)
}

func (c *CategoryUseCase) GetAllCategories() ([]entities.Category, error) {
	return c.categoryRepo.FindAll()
}
