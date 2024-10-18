package repository

import "marketplace/internal/domain/entities"

type CategoryRepository interface {
	Save(category entities.Category, uid int64) error
	FindByID(id int64) (entities.Category, error)
	Update(category entities.Category, uid int64) error
	Delete(id int64, uid int64) error
	FindAllByStore(id int64) ([]entities.Category, error)
}
