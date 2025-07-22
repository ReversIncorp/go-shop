package repository

import (
	"marketplace/internal/domain/entities"
)

type UserRepository interface {
	Create(user *entities.User) (uint64, error)
	FindByEmail(email string) (entities.User, error)
	FindByID(email uint64) (entities.User, error)
}
