package usecase

import (
	"errors"
	"marketplace/internal/domain/models"
	"marketplace/internal/domain/repository"
)

type userUseCaseImpl struct {
	userRepo repository.UserRepository
}

// NewUserUseCase Конструктор для создания новой UserUseCase
func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCaseImpl{userRepo: userRepo}
}

// Register Реализация метода Register
func (u *userUseCaseImpl) Register(user models.User) error {
	existingUser, err := u.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser.ID != "" {
		return errors.New("user already exists")
	}
	return u.userRepo.Create(user)
}

// Login Реализация метода Login
func (u *userUseCaseImpl) Login(email, password string) (string, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil || user.Password != password { // Здесь должна быть логика хэширования пароля
		return "", errors.New("invalid credentials")
	}
	// Генерация токена (здесь может быть использован JWT)
	return "token_placeholder", nil
}

// GetUserByID Реализация метода GetUserByID
func (u *userUseCaseImpl) GetUserByID(id string) (models.User, error) {
	return u.userRepo.FindByID(id)
}
