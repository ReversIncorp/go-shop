package repository

import (
	"marketplace/internal/domain/models"
)

type UserRepository interface {
	Create(user models.User) error
	FindByEmail(email string) (models.User, error)
	FindByID(email string) (models.User, error)
}
