package repository

import (
	"marketplace/internal/domain/entities"
)

type StoreRepository interface {
	Save(store entities.Store, userID int64) error
	FindByID(id int64) (entities.Store, error)
	Update(store entities.Store, userID int64) error
	Delete(id int64, uid int64) error
	FindAll() ([]entities.Store, error)
}
