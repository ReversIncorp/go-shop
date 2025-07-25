package repository

import (
	"errors"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
	errorHandling "marketplace/pkg/error_handling"
	"time"

	"github.com/ztrue/tracerr"
	"gorm.io/gorm"
)

type categoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &categoryRepositoryImpl{
		db: db,
	}
}

func (r *categoryRepositoryImpl) IsExist(id uint64) (bool, error) {
	var category entities.Category
	err := r.db.First(&category, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, tracerr.Wrap(err)
	}
	return true, nil
}

func (r *categoryRepositoryImpl) Save(category *entities.Category) error {
	category.CreatedAt = time.Now()
	if err := r.db.Create(category).Error; err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func (r *categoryRepositoryImpl) Delete(id uint64) error {
	if err := r.db.Delete(&entities.Category{}, id).Error; err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

func (r *categoryRepositoryImpl) FindAllByStore(storeID uint64) ([]entities.Category, error) {
	var store entities.Store
	if err := r.db.Preload("Categories").First(&store, storeID).Error; err != nil {
		return nil, tracerr.Wrap(err)
	}
	return store.Categories, nil
}

func (r *categoryRepositoryImpl) FindByID(id uint64) (entities.Category, error) {
	var category entities.Category
	err := r.db.
		First(&category, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.Category{}, errorHandling.ErrCategoryNotFound
	}
	if err != nil {
		return entities.Category{}, tracerr.Wrap(err)
	}
	return category, nil
}

func (r *categoryRepositoryImpl) FindAll() ([]entities.Category, error) {
	var categories []entities.Category
	if err := r.db.
		Find(&categories).Error; err != nil {
		return nil, tracerr.Wrap(err)
	}
	return categories, nil
}
