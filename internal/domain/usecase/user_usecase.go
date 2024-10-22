package usecase

import (
	"errors"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
	"marketplace/pkg/utils"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (u *UserUseCase) Register(user entities.User) (*entities.Tokens, error) {
	existingUser, err := u.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser.ID != 0 {
		return nil, errors.New("user already exists")
	}

	// Сохраняем пользователя в репозиторий
	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Генерация токенов
	tokenDetails, err := utils.GenerateTokens(user.ID)
	if err != nil {
		return nil, err
	}

	return tokenDetails.ToTokens(), nil
}

func (u *UserUseCase) Login(email, password string) (*entities.Tokens, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil || user.Password != password { // Здесь должна быть логика хэширования пароля
		return nil, errors.New("invalid credentials")
	}

	// Генерация токенов
	tokenDetails, err := utils.GenerateTokens(user.ID)
	if err != nil {
		return nil, err
	}

	return tokenDetails.ToTokens(), nil
}

func (u *UserUseCase) GetUserByID(id uint64) (entities.User, error) {
	return u.userRepo.FindByID(id)
}
