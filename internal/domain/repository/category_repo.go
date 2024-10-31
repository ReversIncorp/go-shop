package repository

import "marketplace/internal/domain/entities"

type CategoryRepository interface {
	Save(category entities.Category) error
	Update(category entities.Category) error
	Delete(id uint64) error
	IsExist(id uint64) (bool, error)

	FindAllByStore(id uint64) ([]entities.Category, error)
}
