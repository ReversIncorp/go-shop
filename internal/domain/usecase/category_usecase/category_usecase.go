package categoryusecase

import (
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
	errorHandling "marketplace/pkg/error_handling"

	"github.com/ztrue/tracerr"
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

func (c *CategoryUseCase) CreateCategory(category *entities.Category, uid uint64) error {
	userData, err := c.userRepo.FindByID(uid)
	if err != nil || !userData.IsSeller {
		return errorHandling.ErrUserNotSeller
	}

	err = c.categoryRepo.Save(category)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func (c *CategoryUseCase) GetCategoryByID(id uint64) (entities.Category, error) {
	category, err := c.categoryRepo.FindByID(id)
	if err != nil {
		return category, tracerr.Wrap(err)
	}
	return category, nil
}

func (c *CategoryUseCase) DeleteCategory(id, uid uint64) error {
	userData, err := c.userRepo.FindByID(uid)
	if err != nil || !userData.IsSeller {
		return errorHandling.ErrUserNotSeller
	}

	exists, err := c.categoryRepo.IsExist(id)
	if err != nil {
		return tracerr.Wrap(err)
	}
	if !exists {
		return errorHandling.ErrCategoryNotFound
	}

	err = c.categoryRepo.Delete(id)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func (c *CategoryUseCase) GetAllCategoriesByStore(id uint64) ([]entities.Category, error) {
	categories, err := c.categoryRepo.FindAllByStore(id)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return categories, nil
}

func (c *CategoryUseCase) GetAllCategories() ([]entities.Category, error) {
	categories, err := c.categoryRepo.FindAll()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return categories, nil
}
