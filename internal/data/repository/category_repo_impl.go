package repository

import (
	"errors"
	"gorm.io/gorm"
	"marketplace/internal/domain/entities"
	repository2 "marketplace/internal/domain/repository"
	"sync"
	"time"
)

type categoryRepositoryImpl struct {
	db *gorm.DB   // Подключение к базе данных
	mu sync.Mutex // Мьютекс для потокобезопасности
}

func NewCategoryRepository(db *gorm.DB) repository2.CategoryRepository {
	return &categoryRepositoryImpl{
		db: db,
	}
}

func (r *categoryRepositoryImpl) Save(category entities.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var existingCategory entities.Category
	if err := r.db.Where("id = ?", category.ID).First(&existingCategory).Error; err == nil {
		return errors.New("category already exists")
	}

	category.CreatedAt = time.Now()
	category.UpdatedAt = category.CreatedAt

	if err := r.db.Create(&category).Error; err != nil {
		return err
	}

	return nil
}

func (r *categoryRepositoryImpl) FindByID(id uint64) (entities.Category, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var category entities.Category
	if err := r.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Category{}, errors.New("category not found")
		}
		return entities.Category{}, err
	}

	return category, nil
}

func (r *categoryRepositoryImpl) Update(category entities.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var existingCategory entities.Category
	if err := r.db.First(&existingCategory, category.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}

	category.UpdatedAt = time.Now()

	if err := r.db.Save(&category).Error; err != nil {
		return err
	}

	return nil
}

func (r *categoryRepositoryImpl) Delete(id uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.db.Delete(&entities.Category{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}

	return nil
}

func (r *categoryRepositoryImpl) FindAll() ([]entities.Category, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var categories []entities.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}
