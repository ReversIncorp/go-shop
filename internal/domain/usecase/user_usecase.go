package usecase

import (
	"marketplace/internal/domain/models"
)

type UserUseCase interface {
	Register(user models.User) error
	Login(email, password string) (string, error) // Возвращает токен
	GetUserByID(id string) (models.User, error)   // Получить пользователя по ID
}
