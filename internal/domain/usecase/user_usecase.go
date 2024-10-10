package usecase

import (
	"errors"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

// NewUserUseCase Конструктор для создания новой UserUseCase
func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

// Register Реализация метода Register
func (u *UserUseCase) Register(user entities.User) error {
	existingUser, err := u.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser.ID != 0 {
		return errors.New("user already exists")
	}
	return u.userRepo.Create(user)
}

// Login Реализация метода Login
func (u *UserUseCase) Login(email, password string) (string, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil || user.Password != password { // Здесь должна быть логика хэширования пароля
		return "", errors.New("invalid credentials")
	}
	// Генерация токена (здесь может быть использован JWT)
	return "token_placeholder", nil
}

// GetUserByID Реализация метода GetUserByID
func (u *UserUseCase) GetUserByID(id uint64) (entities.User, error) {
	return u.userRepo.FindByID(id)
}
