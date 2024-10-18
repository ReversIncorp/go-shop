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

func (c *CategoryUseCase) CreateCategory(category entities.Category, uid int64) error {
	return c.categoryRepo.Save(category, uid)
}

func (c *CategoryUseCase) GetCategoryByID(id int64) (entities.Category, error) {
	return c.categoryRepo.FindByID(id)
}

func (c *CategoryUseCase) UpdateCategory(category entities.Category, uid int64) error {
	return c.categoryRepo.Update(category, uid)
}

func (c *CategoryUseCase) DeleteCategory(id int64, uid int64) error {
	return c.categoryRepo.Delete(id, uid)
}

func (c *CategoryUseCase) GetAllCategoriesByStore(id int64) ([]entities.Category, error) {
	return c.categoryRepo.FindAllByStore(id)
}
