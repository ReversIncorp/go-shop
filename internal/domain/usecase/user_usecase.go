package usecase

import (
	"errors"
	"github.com/labstack/echo/v4"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/enums"
	"marketplace/internal/domain/repository"
	"marketplace/pkg/utils"
)

type UserUseCase struct {
	userRepo  repository.UserRepository
	tokenRepo repository.JWTRepository
}

// NewUserUseCase Конструктор для создания новой UserUseCase
func NewUserUseCase(userRepo repository.UserRepository, tokenRepo repository.JWTRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo, tokenRepo: tokenRepo}
}

// Register Реализация метода Register
func (u *UserUseCase) Register(user entities.User, ctx echo.Context) (*entities.Tokens, error) {
	existingUser, err := u.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser.ID != 0 {
		return nil, errors.New("user already exists")
	}

	// Сохраняем пользователя в репозиторий
	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}
	tokens, err := u.createTokens(user.ID, ctx)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// Login Реализация метода Login
func (u *UserUseCase) Login(email, password string, ctx echo.Context) (*entities.Tokens, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil || user.Password != password { // Здесь должна быть логика хэширования пароля
		return nil, errors.New("invalid credentials")
	}

	tokens, err := u.createTokens(user.ID, ctx)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// GetUserByID Реализация метода GetUserByID
func (u *UserUseCase) GetUserByID(id uint64) (entities.User, error) {
	return u.userRepo.FindByID(id)
}

// UpdateToken Реализация метода GetUserByID
func (u *UserUseCase) UpdateToken(id uint64) (entities.User, error) {
	return u.userRepo.FindByID(id)
}

func (u *UserUseCase) createTokens(userId uint64, ctx echo.Context) (*entities.Tokens, error) {
	{
		accessToken, err := utils.GenerateToken(userId, enums.Access)
		refreshToken, err := utils.GenerateToken(userId, enums.Refresh)
		if err != nil {
			return nil, err
		}

		if err = u.tokenRepo.SaveToken(
			userId,
			accessToken,
			enums.Access,
			ctx,
		); err != nil {
			return nil, err
		}
		if err = u.tokenRepo.SaveToken(
			userId,
			refreshToken,
			enums.Refresh,
			ctx,
		); err != nil {
			return nil, err
		}
		return &entities.Tokens{RefreshToken: refreshToken, AccessToken: accessToken}, nil
	}
}
